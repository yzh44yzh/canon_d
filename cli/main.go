package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"com.github.yzh44yzh/canon_d/internal"
)

// https://en.wikipedia.org/wiki/ANSI_escape_code
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

var reader = bufio.NewReaderSize(os.Stdin, 1<<20)

func main() {
	var deckFile string
	var limit int
	flag.StringVar(&deckFile, "deck", "", "Path to deck file")
	flag.IntVar(&limit, "limit", 10, "Limit the number of cards to learn")
	flag.Parse()

	content, err := os.ReadFile(deckFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO validate limit > 0

	deck := internal.MakeDeck(filepath.Base(deckFile), string(content))
	total := min(limit, len(deck.Cards))
	fmt.Printf("\nDeck %s\"%s\"%s, %d cards\n\n", Blue, deck.Header, Reset, total)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck.Cards), func(i, j int) {
		deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
	})
	deck.Cards = deck.Cards[:limit]

	numErrors := 0

	for i, card := range deck.Cards {
		numErrors += learnCard(card, i+1, total)
	}

	fmt.Printf("%sResults:%s\n  Errors: %d\n", Blue, Reset, numErrors)
}

func learnCard(card internal.Card, idx, total int) int {
	fmt.Printf("%sCard %d/%d%s\n", Yellow, idx, total, Reset)
	fmt.Println(Blue + card.Header.Show() + Reset)

	numErrors := 0
	for _, line := range card.Lines {
		for _, part := range line.Parts {
			if part.Visible {
				fmt.Println(part.Content)
			} else {
				fmt.Print("> ")
				input, _ := reader.ReadString('\n')

				if strings.TrimSpace(input) == strings.TrimSpace(part.Content) {
					fmt.Println(Green + "  OK" + Reset)
				} else {
					numErrors += 1
					fmt.Println(Red + "  Error" + Reset)
					fmt.Println(" ", part.Content)
				}
			}
		}
	}
	_, _ = reader.ReadString('\n')
	return numErrors
}

func clearScreen() {
	homeSeq := "\033[H"
	clearSeq := "\033[2J"
	fmt.Print(homeSeq, clearSeq)
}
