package main

import (
	"fmt"

	"com.github.yzh44yzh/canon_d"
)

func main() {
	card := exampleCard()

	fmt.Println("CARD:")
	fmt.Println(card.Header)
	for _, line := range card.Lines {
		fmt.Println("O:", line.Original)
		fmt.Println("S:", line.Show())
	}
}

func exampleCard() canon_d.Card {
	content := " \n" +
		"### MakeCard\n" +
		"Семейство движков `MergeTree` позволяют схлопывать данные в фоновом режиме.\n" +
		"   \n" +
		"Их 7 штук, из них 3 позволяют реализовать update:\n" +
		"   \n" +
		"- `ReplacingMergeTree`\n" +
		"- `CollapsingMergeTree`\n" +
		"- `VersionedCollapsingMergeTree`\n"

	return canon_d.MakeCard(content)
}
