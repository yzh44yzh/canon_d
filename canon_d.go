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

func MakeLine(content string) Line {
	parts := strings.Split(content, "`")
	lineParts := make([]LinePart, len(parts))
	visible := true

	for i, part := range parts {
		lineParts[i] = LinePart{Visible: visible, Content: part}
		visible = !visible
	}

	return Line{
		Original: content,
		Parts:    lineParts,
	}
}

func (l Line) Show() string {
	var res strings.Builder
	for _, part := range l.Parts {
		if part.Visible {
			res.WriteString(part.Content)
		} else {
			res.WriteString("???")
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
	// TODO MakeCard
	line1 := MakeLine("Семейство движков `MergeTree` позволяют схлопывать данные в фоновом режиме.")
	line2 := MakeLine("Их 7 штук, из них 3 позволяют реализовать update:")
	line3 := MakeLine("- `ReplacingMergeTree`")
	line4 := MakeLine("- `CollapsingMergeTree`")
	line5 := MakeLine("- `VersionedCollapsingMergeTree`")

	return Card{
		Header: "MergeTree",
		Lines:  []Line{line1, line2, line3, line4, line5},
	}
}
