package canon_d

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Deck struct {
	Header string
	Cards  []Card
}

var reBold *regexp.Regexp = regexp.MustCompile(`\*\*(\S.*?)\*\*`) // **bold text**, but __bold text__ is not supported
var reItalic *regexp.Regexp = regexp.MustCompile(`\*(\S.*?)\*`) // *italic text*, but _italic text_ is not supported
var reCode *regexp.Regexp = regexp.MustCompile("`(\\S.*?)`") // `inline code`

type Card struct {
	Header string
	Lines  []Line
}

func MakeCard(content string) Card {
	rawLines := strings.Split(content, "\n")
	header := ""
	lines := []Line{}

	for _, l := range rawLines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		if header == "" {
			header = strings.Trim(l, "#")
			header = strings.TrimSpace(header)
		} else {
			line := MakeLine(l)
			lines = append(lines, line)
		}
	}

	return Card{
		Header: header,
		Lines:  lines,
	}
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
	content := " \n" + 
	"### MakeCard\n" +  
	"Семейство движков `MergeTree` позволяют схлопывать данные в фоновом режиме.\n" +
	"   \n" +
	"Их 7 штук, из них 3 позволяют реализовать update:\n" +
	"   \n" +
	"- `ReplacingMergeTree`\n" +
	"- `CollapsingMergeTree`\n" +
	"- `VersionedCollapsingMergeTree`\n"

	return MakeCard(content)
}
