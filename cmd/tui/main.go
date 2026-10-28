package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"com.github.yzh44yzh/canon_d/internal/parser"
	"com.github.yzh44yzh/canon_d/internal/session"
	"com.github.yzh44yzh/canon_d/internal/tui"
)

func main() {
	deckFile := flag.String("deck", "", "Path to deck file")
	limit := flag.Int("limit", 10, "Limit the number of cards to learn")
	flag.Parse()

	err := run(*deckFile, *limit)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(deckFile string, limit int) error {
	if deckFile == "" {
		return fmt.Errorf("provide a deck with --deck PATH")
	}

	if limit <= 0 {
		return fmt.Errorf("limit must be greater than zero")
	}

	content, err := os.ReadFile(deckFile)
	if err != nil {
		return err
	}

	deck := parser.MakeDeck(filepath.Base(deckFile), string(content))
	session, err := session.NewSession(deck, limit)
	if err != nil {
		return err
	}

	return tui.Run(deck.Header, session)
}
