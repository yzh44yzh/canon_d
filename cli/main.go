package main

import (
	"com.github.yzh44yzh/canon_d"
	"fmt"
)

func main() {
	// deck := canon_d.ExampleDeck()
	// card := deck.Cards[0]
	card := canon_d.ExampleCard()

	fmt.Println(card.Header)
	for _, line := range card.Lines {
		fmt.Println(line.Show())
	}
}
