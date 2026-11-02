// Package decks loads snapshots of Markdown decks from disk.
package decks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"com.github.yzh44yzh/canon_d/internal/model"
	"com.github.yzh44yzh/canon_d/internal/parser"
)

// Entry retains unusable decks so the picker can explain why they are disabled.
type Entry struct {
	Deck model.Deck
	Err  error
}

// Reason omits filesystem paths so disabled rows have room for the actual cause.
func (e Entry) Reason() string {
	if e.Err == nil {
		return ""
	}

	var pathError *os.PathError
	if errors.As(e.Err, &pathError) {
		return pathError.Err.Error()
	}

	return e.Err.Error()
}

func LoadFile(path string) Entry {
	entry := Entry{Deck: model.Deck{Header: filepath.Base(path)}}
	content, err := os.ReadFile(path)
	if err != nil {
		entry.Err = err

		return entry
	}

	entry.Deck = parser.MakeDeck(entry.Deck.Header, string(content))
	if len(entry.Deck.Cards) == 0 {
		entry.Err = errors.New("deck contains no cards")
	}

	return entry
}

// LoadDirectory reads sorted, non-hidden regular .md files without recursion.
// All returned entries are snapshots; starting practice does not reread files.
func LoadDirectory(path string) ([]Entry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var entries []Entry
	var failures []error
	usable := 0

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") || filepath.Ext(file.Name()) != ".md" {
			continue
		}

		if !file.Type().IsRegular() {
			continue
		}

		entry := LoadFile(filepath.Join(path, file.Name()))
		entries = append(entries, entry)
		if entry.Err != nil {
			failures = append(failures, fmt.Errorf("%s: %w", file.Name(), entry.Err))
		} else {
			usable++
		}
	}

	if usable == 0 {
		reason := fmt.Errorf("%s: no usable decks found", path)
		failures = append([]error{reason}, failures...)

		return nil, errors.Join(failures...)
	}

	return entries, nil
}
