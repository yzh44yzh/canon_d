package history

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"com.github.yzh44yzh/canon_d/internal/model"
	"github.com/google/go-cmp/cmp"
)

func openStore(t *testing.T, dir string) *Store {
	t.Helper()

	store, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { store.Close() })

	return store
}

func TestSaveAcrossRunsAndCheckpoints(t *testing.T) {
	dir := t.TempDir()
	id := model.CardIdentity{Deck: "a,\"b.md", Header: "Parent / Child, \"quoted\""}
	reading := model.CardIdentity{Deck: "reading.md", Header: "Read"}
	store, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if err := store.Save(nil); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(store.path); !os.IsNotExist(err) {
		t.Fatal("quitting without showing cards must not create results")
	}

	first := statistics{
		id:      {Shows: 1, CorrectAnswers: 2, IncorrectAnswers: 1},
		reading: {Shows: 1},
	}
	if err := store.Save(first); err != nil {
		t.Fatal(err)
	}

	// The retry checkpoint includes all results since startup, not just its delta.
	first[id] = model.CardStatistics{Shows: 2, CorrectAnswers: 5, IncorrectAnswers: 1}
	for i := 0; i < 2; i++ {
		if err := store.Save(first); err != nil {
			t.Fatal(err)
		}
	}

	if err := store.Close(); err != nil {
		t.Fatal(err)
	}

	store = openStore(t, dir)
	second := statistics{id: {Shows: 1, IncorrectAnswers: 1}}
	if err := store.Save(second); err != nil {
		t.Fatal(err)
	}

	got, err := readCSV(store.path)
	if err != nil {
		t.Fatal(err)
	}

	want := statistics{
		id:      {Shows: 3, CorrectAnswers: 5, IncorrectAnswers: 2},
		reading: {Shows: 1},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("wrong persisted totals (-want +got):\n%s", diff)
	}

	matches, err := filepath.Glob(filepath.Join(dir, ".canon_d-*.tmp"))
	if err != nil || len(matches) != 0 {
		t.Fatalf("temporary files were not cleaned up: %v, %v", matches, err)
	}
}

func TestExclusiveFolderLock(t *testing.T) {
	dir := t.TempDir()
	first, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()

	if second, err := Open(dir); err == nil {
		second.Close()
		t.Fatal("a second writer must be refused")
	}

	other := openStore(t, t.TempDir())
	if err := other.Save(nil); err != nil {
		t.Fatal(err)
	}

	if err := first.Close(); err != nil {
		t.Fatal(err)
	}

	openStore(t, dir)
}

func TestInitialStatisticsAreAnIndependentStartupSnapshot(t *testing.T) {
	dir := t.TempDir()
	id := model.CardIdentity{Deck: "a.md", Header: "Card"}
	want := statistics{id: {Shows: 3, CorrectAnswers: 1, IncorrectAnswers: 2}}
	if err := writeCSV(filepath.Join(dir, ".canon_d.csv"), want); err != nil {
		t.Fatal(err)
	}

	store := openStore(t, dir)
	snapshot := store.InitialStatistics()
	delete(snapshot, id)
	if err := store.Save(statistics{id: {Shows: 1, CorrectAnswers: 1}}); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, store.InitialStatistics()); diff != "" {
		t.Fatalf("startup statistics changed after a save or caller mutation: %s", diff)
	}
}

func TestInvalidCSVPreservedAndLockReleased(t *testing.T) {
	tests := map[string]string{
		"empty":      "",
		"header":     "wrong,header\n",
		"negative":   "a.md,Card,1,-1,0\n",
		"not number": "a.md,Card,abc,0,0\n",
		"overflow":   "a.md,Card,18446744073709551616,0,0\n",
		"short row":  "a.md,Card,1,0\n",
		"duplicate":  "a.md,Card,1,0,0\na.md,Card,1,0,0\n",
		"quoting":    "a.md,\"Card,1,0,0\n",
	}

	for name, rows := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, ".canon_d.csv")
			content := rows
			if name != "empty" && name != "header" {
				content = strings.Join(columns, ",") + "\n" + rows
			}

			if err := os.WriteFile(path, []byte(content), 0600); err != nil {
				t.Fatal(err)
			}

			if store, err := Open(dir); err == nil {
				store.Close()
				t.Fatal("invalid CSV must prevent startup")
			}

			got, err := os.ReadFile(path)
			if err != nil || string(got) != content {
				t.Fatal("invalid CSV must be preserved")
			}

			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}

			openStore(t, dir)
		})
	}
}

func TestFailedSavePreservesPreviousFile(t *testing.T) {
	dir := t.TempDir()
	store := openStore(t, dir)
	id := model.CardIdentity{Deck: "a.md", Header: "Card"}
	initial := statistics{id: {Shows: math.MaxUint64}}
	if err := store.Save(initial); err != nil {
		t.Fatal(err)
	}

	before, err := os.ReadFile(store.path)
	if err != nil {
		t.Fatal(err)
	}

	// A later run cannot add to a counter that has reached its maximum.
	store.base = initial
	if err := store.Save(statistics{id: {Shows: 1}}); err == nil {
		t.Fatal("overflow must fail rather than wrap")
	}

	// A failed replacement leaves the previous file available and removes the temp.
	path := store.path
	store.path = filepath.Join(dir, "missing", ".canon_d.csv")
	store.base = make(statistics)
	if err := store.Save(statistics{id: {Shows: 1}}); err == nil {
		t.Fatal("unwritable destination must fail")
	}

	after, err := os.ReadFile(path)
	if err != nil || string(before) != string(after) {
		t.Fatal("failed save changed previously saved data")
	}
}
