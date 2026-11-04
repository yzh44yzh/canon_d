package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"com.github.yzh44yzh/canon_d/internal/decks"
	"com.github.yzh44yzh/canon_d/internal/history"
	"com.github.yzh44yzh/canon_d/internal/parser"
	"com.github.yzh44yzh/canon_d/internal/session"
	"com.github.yzh44yzh/canon_d/internal/tui"
)

func main() {
	os.Exit(runCLI(os.Args[1:], os.Stdout, os.Stderr))
}

func runCLI(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("canon_d", flag.ContinueOnError)
	flags.SetOutput(stderr)
	deckFile := flags.String("deck", "", "Path to deck file")
	deckDir := flags.String("decks", "", "Path to folder containing deck files")
	limit := flags.Int("limit", 10, "Limit the number of cards to learn")
	validateFile := flags.String("validate-deck", "",
		"Validate a deck without starting the TUI")

	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}

		return 2
	}

	provided := make(map[string]bool)
	flags.Visit(func(f *flag.Flag) { provided[f.Name] = true })

	if flags.NArg() > 0 {
		fmt.Fprintln(stderr, "positional arguments are not supported")
		return 2
	}

	if provided["validate-deck"] {
		conflict := provided["deck"] || provided["decks"] || provided["limit"]
		if *validateFile == "" || conflict {
			fmt.Fprintln(stderr, "use --validate-deck FILE alone")
			return 2
		}

		return validateDeck(*validateFile, stdout, stderr)
	}

	if err := validateInput(*deckFile, *deckDir, *limit); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	err := run(*deckFile, *deckDir, *limit)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	return 0
}

func validateDeck(path string, stdout, stderr io.Writer) int {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	report := parser.Validate(string(content))

	for _, finding := range report.Findings {
		severity := "error"
		if finding.Warning {
			severity = "warning"
		}

		fmt.Fprintf(stdout, "%s:%d: %s: %s", path, finding.Line, severity, finding.Message)

		if finding.RelatedLine > 0 {
			fmt.Fprintf(stdout, " (first occurrence: %s:%d)", path, finding.RelatedLine)
		}

		fmt.Fprintln(stdout)
	}

	fmt.Fprintf(stdout, "%s: %d cards, %d reading cards, %d warnings, %d errors\n",
		path, report.Cards, report.ReadingCards, report.Warnings, report.Errors)

	if report.Errors > 0 {
		return 1
	}

	return 0
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

		return tui.RunSelection(entries, limit, store.InitialStatistics(), save)
	}

	entry := decks.LoadFile(deckFile)
	if entry.Err != nil {
		return entry.Err
	}

	session, err := session.NewSession(entry.Deck, limit, store.InitialStatistics())
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
