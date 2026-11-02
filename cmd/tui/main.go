package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"com.github.yzh44yzh/canon_d/internal/decks"
	"com.github.yzh44yzh/canon_d/internal/history"
	"com.github.yzh44yzh/canon_d/internal/session"
	"com.github.yzh44yzh/canon_d/internal/tui"
)

func main() {
	deckFile := flag.String("deck", "", "Path to deck file")
	deckDir := flag.String("decks", "", "Path to folder containing deck files")
	limit := flag.Int("limit", 10, "Limit the number of cards to learn")
	flag.Parse()

	err := run(*deckFile, *deckDir, *limit)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(deckFile, deckDir string, limit int) error {
	if err := validateInput(deckFile, deckDir, limit); err != nil {
		return err
	}

	folder := deckDir
	if folder == "" {
		folder = filepath.Dir(deckFile)
	}

	store, err := history.Open(folder)
	if err != nil {
		return err
	}
	defer store.Close()

	save := func(s *session.Session) error { return store.Save(s.Statistics()) }

	if deckDir != "" {
		entries, err := decks.LoadDirectory(deckDir)
		if err != nil {
			return err
		}

		return tui.RunSelection(entries, limit, save)
	}

	entry := decks.LoadFile(deckFile)
	if entry.Err != nil {
		return entry.Err
	}

	session, err := session.NewSession(entry.Deck, limit)
	if err != nil {
		return err
	}

	return tui.Run(session, save)
}

func validateInput(deckFile, deckDir string, limit int) error {
	if deckFile == "" && deckDir == "" {
		return fmt.Errorf("provide --deck FILE or --decks DIR")
	}

	if deckFile != "" && deckDir != "" {
		return fmt.Errorf("use either --deck FILE or --decks DIR, not both")
	}

	if limit <= 0 {
		return fmt.Errorf("limit must be greater than zero")
	}

	return nil
}
