package internal

import (
	"errors"
	"regexp"
	"strings"
)

func MakeDeck(header, content string) Deck {
	var currHeader *CardHeader = nil
	cards := []Card{}
	cardLines := []string{}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		// end of the current card,
		if strings.HasPrefix(line, "#") {
			card, err := makeCard(currHeader, cardLines)
			if err == nil {
				cards = append(cards, card)
			}
			cardLines = cardLines[:0]

			// start for the next card
			header := makeCardHeader(line)
			for {
				if currHeader == nil {
					currHeader = &header
					break
				}
				if currHeader.Level < header.Level {
					header.Parent = currHeader
					currHeader = &header
					break
				} else {
					currHeader = currHeader.Parent
				}
			}
		} else {
			cardLines = append(cardLines, line)
		}
	}

	// last card
	card, err := makeCard(currHeader, cardLines)
	if err == nil {
		cards = append(cards, card)
	}

	return Deck{
		Header: header,
		Cards:  cards,
	}
}

func makeCardHeader(line string) CardHeader {
	level := 0
	for strings.HasPrefix(line, "#") {
		level += 1
		line = line[1:]
	}

	return CardHeader{
		Level:   level,
		Content: strings.TrimSpace(line),
	}
}

func makeCard(header *CardHeader, rawLines []string) (Card, error) {
	card := Card{}

	if header == nil {
		return card, errors.New("no header")
	}

	card.Header = *header

	var multiLineCode bool = false
	var line Line

	for _, l := range rawLines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		if strings.HasPrefix(l, "```") {
			multiLineCode = !multiLineCode
			continue
		}

		if multiLineCode {
			line = makeInvisibleLine(l)
		} else {
			line = makeLine(l)
		}
		card.Lines = append(card.Lines, line)
	}

	if len(card.Lines) == 0 {
		return card, errors.New("empty card")
	}
	return card, nil
}

func makeLine(content string) Line {
	lineParts := []LinePart{}
	rest := content

	for {
		before, part, after, found := findPart(rest)
		if found {
			if before != "" {
				beforePart := LinePart{
					Content: before,
					Visible: true,
				}
				lineParts = append(lineParts, beforePart)
			}

			currPart := LinePart{
				Content: part,
				Visible: false,
			}
			lineParts = append(lineParts, currPart)

			rest = after
		} else {
			if rest != "" {
				restPart := LinePart{
					Content: rest,
					Visible: true,
				}
				lineParts = append(lineParts, restPart)
			}
			break
		}
	}

	return Line{
		Original: content,
		Parts:    lineParts,
	}
}

func makeInvisibleLine(content string) Line {
	part := LinePart{
		Content: content,
		Visible: false,
	}
	return Line{
		Original: content,
		Parts:    []LinePart{part},
	}
}

// `inline code`
var reCode *regexp.Regexp = regexp.MustCompile("`(\\S.*?)`")

func findPart(content string) (before, part, after string, found bool) {
	idx := reCode.FindStringIndex(content)
	if idx != nil {
		f := idx[0]
		t := idx[1]
		return content[:f], content[f+1 : t-1], content[t:], true
	}

	return "", "", "", false
}
