package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name, file, dir string
		limit           int
		wantError       bool
	}{
		{name: "file", file: "deck.md", limit: 10},
		{name: "folder", dir: "decks", limit: 10},
		{name: "missing", limit: 10, wantError: true},
		{name: "both", file: "deck.md", dir: "decks", limit: 10, wantError: true},
		{name: "zero", dir: "decks", wantError: true},
		{name: "negative", file: "deck.md", limit: -1, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInput(tt.file, tt.dir, tt.limit)
			if (err != nil) != tt.wantError {
				t.Fatalf("unexpected validation result: %v", err)
			}
		})
	}
}

func TestValidationCLI(t *testing.T) {
	tests := []struct {
		name, content, output string
		code                  int
	}{
		{
			name: "valid", content: "# Card\n`answer`",
			output: ": 1 cards, 0 reading cards, 0 warnings, 0 errors\n",
		},
		{
			name: "warnings", content: "intro\n# Card\nread",
			output: ": 1 cards, 1 reading cards, 1 warnings, 0 errors\n",
		},
		{
			name: "errors", content: "# Card\nread\n# Card\nread", code: 1,
			output: ":3: error: duplicate card path: Card (first occurrence: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "deck.md")

			if err := os.WriteFile(path, []byte(tt.content), 0600); err != nil {
				t.Fatal(err)
			}

			// A corrupt history file must not prevent validation or be rewritten.
			historyPath := filepath.Join(dir, ".canon_d.csv")
			err := os.WriteFile(historyPath, []byte("invalid history"), 0600)
			if err != nil {
				t.Fatal(err)
			}

			var stdout, stderr bytes.Buffer
			code := runCLI([]string{"--validate-deck", path}, &stdout, &stderr)

			if code != tt.code || stderr.Len() != 0 {
				t.Fatalf("exit %d, stderr %q; want exit %d", code, &stderr, tt.code)
			}

			if !strings.Contains(stdout.String(), path+tt.output) {
				t.Fatalf("unexpected stdout: %s", &stdout)
			}

			if tt.code == 1 && !strings.Contains(stdout.String(), path+":1)") {
				t.Fatalf("missing first duplicate location: %s", &stdout)
			}

			entries, err := os.ReadDir(dir)
			if err != nil || len(entries) != 2 {
				t.Fatalf("validation created files: %v, %v", entries, err)
			}

			data, err := os.ReadFile(historyPath)
			if err != nil || string(data) != "invalid history" {
				t.Fatalf("validation changed history: %q, %v", data, err)
			}
		})
	}
}

func TestValidationCLIUsage(t *testing.T) {
	tests := [][]string{
		{"--validate-deck"},
		{"--validate-deck="},
		{"--validate-deck", "missing.md", "--deck", "other.md"},
		{"--validate-deck", "missing.md", "--deck="},
		{"--validate-deck", "missing.md", "--decks", "decks"},
		{"--validate-deck", "missing.md", "--limit", "10"},
		{"--limit", "0", "--validate-deck", "missing.md"},
		{"--validate-deck", "missing.md", "extra"},
		{"--validate-deck", "missing.md"},
		{"--unknown"},
	}

	for _, args := range tests {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := runCLI(args, &stdout, &stderr); code != 2 {
				t.Fatalf("exit %d; want 2", code)
			}

			if stdout.Len() != 0 || stderr.Len() == 0 {
				t.Fatalf("stdout %q, stderr %q", &stdout, &stderr)
			}
		})
	}
}

func TestInvalidHistoryStopsPracticeStartup(t *testing.T) {
	for _, mode := range []string{"--deck", "--decks"} {
		t.Run(mode, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, ".canon_d.csv")
			content := "deck,card_header,shows,correct_answers,incorrect_answers\n" +
				"removed.md,Old card,0,1,0\n"
			if err := os.WriteFile(path, []byte(content), 0600); err != nil {
				t.Fatal(err)
			}

			target := dir
			if mode == "--deck" {
				target = filepath.Join(dir, "missing.md")
			}

			// Missing decks must not mask history validation or reach TUI startup.
			var stdout, stderr bytes.Buffer
			code := runCLI([]string{mode, target}, &stdout, &stderr)
			if code != 1 || stdout.Len() != 0 {
				t.Fatalf("exit %d, stdout %q", code, &stdout)
			}

			if !strings.Contains(stderr.String(), path+":2: invalid correct_answers") {
				t.Fatalf("unexpected diagnostic: %s", &stderr)
			}

			got, err := os.ReadFile(path)
			if err != nil || string(got) != content {
				t.Fatal("invalid history was changed")
			}
		})
	}
}
