package decks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func writeDeck(t *testing.T, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestDirectoryDiscoveryAndSnapshot(t *testing.T) {
	dir := t.TempDir()
	writeDeck(t, dir, "b.md", "# B\n`original`")
	writeDeck(t, dir, "a.md", "# A\nReading")
	writeDeck(t, dir, "empty.md", "# Empty")
	writeDeck(t, dir, ".hidden.md", "# Hidden\n`hidden`")
	writeDeck(t, dir, "ignore.txt", "# Ignored\n`ignore`")
	sub := filepath.Join(dir, "sub.md")
	if err := os.Mkdir(sub, 0700); err != nil {
		t.Fatal(err)
	}

	writeDeck(t, sub, "nested.md", "# Nested\n`nested`")
	entries, err := LoadDirectory(dir)
	if err != nil {
		t.Fatal(err)
	}

	var names []string

	for _, entry := range entries {
		names = append(names, entry.Deck.Header)
	}

	if diff := cmp.Diff([]string{"a.md", "b.md", "empty.md"}, names); diff != "" {
		t.Fatalf("unexpected discovered decks (-want +got):\n%s", diff)
	}

	if entries[0].Err != nil || entries[1].Err != nil || entries[2].Err == nil {
		t.Fatal("only the empty deck should be disabled")
	}

	writeDeck(t, dir, "b.md", "# Changed\n`changed`")
	if err := os.Remove(filepath.Join(dir, "a.md")); err != nil {
		t.Fatal(err)
	}

	if entries[1].Deck.Cards[0].Lines[0].Original != "`original`" ||
		len(entries[0].Deck.Cards) != 1 {
		t.Fatal("loaded contents must survive subsequent file changes and removal")
	}
}

func TestLoadingErrors(t *testing.T) {
	dir := t.TempDir()
	if _, err := LoadDirectory(dir); err == nil ||
		!strings.Contains(err.Error(), "no usable decks") {
		t.Fatalf("expected explanatory empty-directory error, got %v", err)
	}

	writeDeck(t, dir, "empty.md", "# Empty")
	if _, err := LoadDirectory(dir); err == nil ||
		!strings.Contains(err.Error(), "empty.md") {
		t.Fatalf("expected disabled-deck reason, got %v", err)
	}

	missing := filepath.Join(dir, "missing.md")
	if entry := LoadFile(missing); entry.Err == nil || entry.Deck.Header != "missing.md" {
		t.Fatalf("expected unreadable deck with its filename, got %+v", entry)
	}

	if _, err := LoadDirectory(missing); err == nil {
		t.Fatal("expected directory read failure")
	}
}
