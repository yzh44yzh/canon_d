package main

import (
	_ "embed"
	"fmt"

	"com.github.yzh44yzh/canon_d"
)

//go:embed example.md
var exampleDeck string

func main() {
	deck := canon_d.MakeDeck("Example", exampleDeck)

	for _, card := range deck.Cards {
		fmt.Println("---")
		fmt.Println(card.Header)
		for _, line := range card.Lines {
			fmt.Println(line.Show())
		}
	}
}
