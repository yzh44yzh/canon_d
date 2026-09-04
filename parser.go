package canon_d

import (
	"errors"
	"regexp"
	"strings"
)

// **bold text**, but __bold text__ is not supported
var reBold *regexp.Regexp = regexp.MustCompile(`\*\*(\S.*?)\*\*`)

// _italic text_, but *italic text* is not supported
var reItalic *regexp.Regexp = regexp.MustCompile(`\_(\S.*?)\_`)

// `inline code`
var reCode *regexp.Regexp = regexp.MustCompile("`(\\S.*?)`")

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

	for _, l := range rawLines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		line := makeLine(l)
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

func findPart(content string) (before, part, after string, found bool) {
	var minIdx []int = nil
	cut := 0

	idxB := reBold.FindStringIndex(content)
	if idxB != nil {
		minIdx = idxB
		cut = 2
	}

	idxI := reItalic.FindStringIndex(content)
	switch {
	case idxI != nil && minIdx == nil:
		minIdx = idxI
		cut = 1
	case idxI != nil && idxI[0] < minIdx[0]:
		minIdx = idxI
		cut = 1
	}

	idxC := reCode.FindStringIndex(content)
	switch {
	case idxC != nil && minIdx == nil:
		minIdx = idxC
		cut = 1
	case idxC != nil && idxC[0] < minIdx[0]:
		minIdx = idxC
		cut = 1
	}

	if minIdx == nil {
		return "", "", "", false
	}

	f := minIdx[0]
	t := minIdx[1]

	return content[:f], content[f+cut : t-cut], content[t:], true
}
