package main

import (
	_ "embed"
	"fmt"

	"com.github.yzh44yzh/canon_d"
)

//go:embed example.md
var exampleDeck string

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

func main() {
	deck := canon_d.MakeDeck("Example", exampleDeck)
	total := len(deck.Cards)
	fmt.Printf("Deck %s\"%s\"%s, %d cards\n", Blue, deck.Header, Reset, total)

	var input string
	fmt.Scanln(&input)
	clearScreen()

	for i, card := range deck.Cards {
		learnCard(card, i+1, total)
	}

	fmt.Printf("%sResults:%s\ntodo\n", Blue, Reset)
}

func learnCard(card canon_d.Card, idx, total int) {
	fmt.Printf("%sCard %d/%d%s\n", Yellow, idx, total, Reset)
	fmt.Println(Blue + card.Header.Show() + Reset)
	for _, line := range card.Lines {
		fmt.Println(line.Show())
	}

	var input string
	fmt.Scanln(&input)
	clearScreen()
}

func clearScreen() {
	homeSeq := "\033[H"
	clearSeq := "\033[2J"
	fmt.Print(homeSeq, clearSeq)
}
