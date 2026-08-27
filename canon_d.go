package canon_d

import (
	"strings"
	"unicode/utf8"
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

type Line struct {
	Original string
	Parts    []LinePart
}

type LinePart struct {
	Visible bool
	Content string
}

func MakeLine(content string) Line {
	lineParts := []LinePart{}
	currPart := LinePart{
		Content: "",
		Visible: true,
	}

	for content != "" {
		after, found := cutSep(content)
		if found {
			if currPart.Content != "" {
				lineParts = append(lineParts, currPart)
			}
			nextPart := LinePart{
				Content: "",
				Visible: !currPart.Visible,
			}
			currPart = nextPart
			content = after
		} else {
			firstRune, _ := utf8.DecodeRuneInString(content)
			first := string(firstRune)
			currPart.Content += first
			content = strings.TrimPrefix(content, first)
		}
	}
	if currPart.Content != "" {
		lineParts = append(lineParts, currPart)
	}

	return Line{
		Original: content,
		Parts:    lineParts,
	}
}

func cutSep(str string) (string, bool) {
	separators := []string{
		"`",  // code
		"**", // bold
		"__", // bold
	}
	for _, sep := range separators {
		after, found := strings.CutPrefix(str, sep)
		if found {
			return after, true
		}
	}

	return "", false
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
