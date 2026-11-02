package main

import "testing"

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
