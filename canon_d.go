package canon_d

import (
	"strings"
)

type Deck struct {
	Header string
	Cards  []Card
}

// TODO Cards group in trees by different header levels

type Card struct {
	Header string
	Lines  []Line
}

// TODO CardPart :: Line | CodeBlock

type Line struct {
	Original string
	Parts    []LinePart
}

type LinePart struct {
	Visible bool
	Content string
}

func (l Line) Show() string {
	var res strings.Builder
	last := len(l.Parts) - 1

	for i, part := range l.Parts {
		if part.Visible {
			res.WriteString(part.Content)
		} else {
			switch {
			case i == last:
				res.WriteString("???")
			case i == 0:
				res.WriteString("??? ")
			default:
				res.WriteString(" ??? ")
			}
		}
	}
	return res.String()
}

func ExampleDeck() Deck {
	card := ExampleCard()

	return Deck{
		Header: "Example",
		Cards:  []Card{card},
	}
}

func ExampleCard() Card {
	// ### MergeTree
	// Семейство движков `MergeTree` позволяют схлопывать данные в фоновом режиме.
	// Их 7 штук, из них 3 позволяют реализовать update.
	// `ReplacingMergeTree`
	// `CollapsingMergeTree`
	// `VersionedCollapsingMergeTree`

	content1 := "Семейство движков"
	linePart1 := LinePart{Visible: true, Content: content1}

	content2 := "MergeTree"
	linePart2 := LinePart{Visible: false, Content: content2}

	content3 := "позволяют схлопывать данные в фоновом режиме."
	linePart3 := LinePart{Visible: true, Content: content3}

	line1 := Line{
		Original: content1 + " " + content2 + " " + content3,
		Parts:    []LinePart{linePart1, linePart2, linePart3},
	}

	content4 := "Их 7 штук, из них 3 позволяют реализовать update."
	linePart4 := LinePart{Visible: true, Content: content4}

	line2 := Line{
		Original: content4,
		Parts:    []LinePart{linePart4},
	}

	content5 := "ReplacingMergeTree"
	linePart5 := LinePart{Visible: false, Content: content5}

	line3 := Line{
		Original: content5,
		Parts:    []LinePart{linePart5},
	}

	content6 := "CollapsingMergeTree"
	linePart6 := LinePart{Visible: false, Content: content6}

	line4 := Line{
		Original: content6,
		Parts:    []LinePart{linePart6},
	}

	content7 := "VersionedCollapsingMergeTree"
	linePart7 := LinePart{Visible: false, Content: content7}

	line5 := Line{
		Original: content7,
		Parts:    []LinePart{linePart7},
	}

	return Card{
		Header: "MergeTree",
		Lines:  []Line{line1, line2, line3, line4, line5},
	}
}
