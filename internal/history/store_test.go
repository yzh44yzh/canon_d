package history

import (
	"fmt"
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

func TestSaveRejectsInvalidMergedRows(t *testing.T) {
	tests := []struct {
		name, value, reason string
		field               int
	}{
		{name: "empty deck", field: 0, reason: "non-whitespace"},
		{name: "blank header", field: 1, value: " \u2003", reason: "non-whitespace"},
		{name: "invalid UTF-8 deck", field: 0, value: "\xff.md", reason: "UTF-8"},
		{name: "invalid UTF-8 header", field: 1, value: "\xff", reason: "UTF-8"},
		{name: "deck newline", field: 0, value: "a\nb", reason: "control"},
		{name: "header newline", field: 1, value: "a\nb", reason: "control"},
		{name: "header tab", field: 1, value: "a\tb", reason: "control"},
		{name: "header NUL", field: 1, value: "a\x00b", reason: "control"},
		{name: "slash", field: 0, value: "a/b", reason: "filename"},
		{name: "backslash", field: 0, value: "a\\b", reason: "filename"},
		{name: "dot", field: 0, value: ".", reason: "filename"},
		{name: "dot dot", field: 0, value: "..", reason: "filename"},
		{
			name: "zero shows correct", field: 3,
			reason: "positive answers require shows > 0",
		},
		{
			name: "zero shows incorrect", field: 4,
			reason: "positive answers require shows > 0",
		},
	}

	for _, existing := range []bool{false, true} {
		for _, tt := range tests {
			t.Run(fmt.Sprintf("existing=%v/%s", existing, tt.name), func(t *testing.T) {
				dir := t.TempDir()
				path := filepath.Join(dir, ".canon_d.csv")
				var before string

				if existing {
					row := []string{"old.md", "Removed", "2", "3", "4"}
					before = "\n" + encodeHistory(t, row)
					err := os.WriteFile(path, []byte(before), 0600)
					if err != nil {
						t.Fatal(err)
					}
				}

				store := openStore(t, dir)
				initial := store.InitialStatistics()
				id := model.CardIdentity{Deck: "test.md", Header: "Card"}
				total := model.CardStatistics{Shows: 1}

				switch tt.field {
				case 0:
					id.Deck = tt.value
				case 1:
					id.Header = tt.value
				case 3:
					total = model.CardStatistics{CorrectAnswers: 1}
				case 4:
					total = model.CardStatistics{IncorrectAnswers: 1}
				}

				results := statistics{
					{Deck: "good.md", Header: "Other"}: {Shows: 1},
					id:                                 total,
				}

				for i := 0; i < 2; i++ {
					err := store.Save(results)
					if err == nil {
						t.Fatal("invalid save accepted or cached")
					}

					identity := fmt.Sprintf("deck %q, card %q",
						id.Deck, id.Header)
					checks := []string{
						path, identity, columns[tt.field], tt.reason,
					}
					for _, text := range checks {
						if !strings.Contains(err.Error(), text) {
							t.Fatalf("error %q lacks %q",
								err, text)
						}
					}
				}

				got, err := os.ReadFile(path)
				if existing && (err != nil || string(got) != before) {
					t.Fatal("invalid save changed existing history bytes")
				}

				if !existing && !os.IsNotExist(err) {
					t.Fatal("invalid first save must not create history")
				}

				matches, err := filepath.Glob(filepath.Join(dir, ".canon_d-*.tmp"))
				if err != nil || len(matches) != 0 {
					t.Fatalf("temporary files: %v, %v", matches, err)
				}

				diff := cmp.Diff(initial, store.InitialStatistics())
				if diff != "" {
					t.Fatalf("startup snapshot changed: %s", diff)
				}

				if err := store.Close(); err != nil {
					t.Fatal(err)
				}

				reopened := openStore(t, dir)
				diff = cmp.Diff(initial, reopened.InitialStatistics())
				if diff != "" {
					t.Fatalf("invalid save changed reopened history: %s", diff)
				}
			})
		}
	}
}

func TestSaveValidatesMergedTotalsAndReopens(t *testing.T) {
	dir := t.TempDir()
	id := model.CardIdentity{Deck: "a", Header: "Card"}
	stale := model.CardIdentity{Deck: "removed.md", Header: "Old / Card"}
	base := statistics{id: {Shows: 1}, stale: {Shows: 2, IncorrectAnswers: 3}}
	path := filepath.Join(dir, ".canon_d.csv")
	if err := writeCSV(path, base); err != nil {
		t.Fatal(err)
	}

	store := openStore(t, dir)
	results := statistics{
		id:                            {CorrectAnswers: 2},
		{Deck: "A", Header: " Card "}: {},
		{Deck: "old,\"deck\"", Header: "Go / Child, \"quoted\""}: {
			Shows: 1, CorrectAnswers: math.MaxUint64, IncorrectAnswers: math.MaxUint64,
		},
	}
	if err := store.Save(results); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(base, store.InitialStatistics()); diff != "" {
		t.Fatalf("valid save changed startup snapshot: %s", diff)
	}

	// Repeated checkpoints still merge against the startup base, not the last save.
	if err := store.Save(results); err != nil {
		t.Fatal(err)
	}

	if err := store.Close(); err != nil {
		t.Fatal(err)
	}

	want := statistics{
		id:                            {Shows: 1, CorrectAnswers: 2},
		stale:                         {Shows: 2, IncorrectAnswers: 3},
		{Deck: "A", Header: " Card "}: {},
		{Deck: "old,\"deck\"", Header: "Go / Child, \"quoted\""}: {
			Shows: 1, CorrectAnswers: math.MaxUint64, IncorrectAnswers: math.MaxUint64,
		},
	}
	reopened := openStore(t, dir)
	if diff := cmp.Diff(want, reopened.InitialStatistics()); diff != "" {
		t.Fatalf("save/reopen mismatch (-want +got):\n%s", diff)
	}
}

func TestSaveReportsFirstInvalidRowInIdentityOrder(t *testing.T) {
	store := openStore(t, t.TempDir())
	results := statistics{
		{Deck: "z.md", Header: "\t"}:    {Shows: 1},
		{Deck: "a.md", Header: "Bad\t"}: {Shows: 1},
		{Deck: "a.md", Header: ""}:      {Shows: 1},
	}
	err := store.Save(results)
	want := fmt.Sprintf("save statistics: %s, deck %q, card %q: ", store.path, "a.md", "") +
		"invalid card_header: must contain non-whitespace text"
	if err == nil || err.Error() != want {
		t.Fatalf("save error = %v; want %s", err, want)
	}
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
