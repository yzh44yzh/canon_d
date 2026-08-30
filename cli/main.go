package main

import (
	"fmt"
	"regexp"

	"com.github.yzh44yzh/canon_d"
)

func main() {
	// deck := canon_d.ExampleDeck()
	// card := deck.Cards[0]
	card := canon_d.ExampleCard()

	fmt.Println(card.Header)
	for _, line := range card.Lines {
		fmt.Println(line.Show())
	}

	tryRegexp()
}

func tryRegexp() {
	content := "one _two_ three _four_ five **six** seven `eight` ten"

	re := regexp.MustCompile(`_(\S.*?)_`)
	res := re.FindString(content)
	fmt.Println(res) // '_two_'

	re2 := regexp.MustCompile(`\*\*(\S.*?)\*\*`)
	res2 := re2.FindString(content)
	fmt.Println(res2) // **six**

	re3 := regexp.MustCompile("`(\\S.*?)`")
	res3 := re3.FindString(content)
	fmt.Println(res3) // `eight`
}
