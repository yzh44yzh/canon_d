package parser

import (
	"regexp"
	"strings"

	"com.github.yzh44yzh/canon_d/internal/model"
)

// MakeDeck parses Markdown cards while preserving their formatting.
func MakeDeck(header, content string) model.Deck {
	content = strings.ReplaceAll(content, "\r\n", "\n")

	deck := model.Deck{Header: header}
	var currHeader *model.CardHeader
	var lines []string
	var fence string

	maybeAddCard := func() {
		if currHeader == nil || strings.TrimSpace(strings.Join(lines, "\n")) == "" {
			return
		}

		newCard := formattedCard(*currHeader, lines)
		deck.Cards = append(deck.Cards, newCard)
	}

	for _, line := range strings.Split(content, "\n") {
		if fence == "" && strings.HasPrefix(line, "#") {
			maybeAddCard()

			lines = nil
			nextHeader := makeCardHeader(line)

			for currHeader != nil && currHeader.Level >= nextHeader.Level {
				currHeader = currHeader.Parent
			}

			nextHeader.Parent = currHeader
			currHeader = &nextHeader
			continue
		}

		lines = append(lines, line)
		fence = nextFence(line, fence)
	}

	maybeAddCard()

	return deck
}

// nextFence recognizes backtick and tilde fences, including longer fences that
// can contain shorter fence markers as answer text.
func nextFence(line, fence string) string {
	trimmed := strings.TrimSpace(line)
	if fence != "" {
		if len(trimmed) >= len(fence) && strings.Trim(trimmed, fence[:1]) == "" {
			return ""
		}

		return fence
	}

	if len(trimmed) < 3 || (trimmed[0] != '`' && trimmed[0] != '~') {
		return ""
	}

	n := 0

	for n < len(trimmed) && trimmed[n] == trimmed[0] {
		n++
	}

	if n < 3 {
		return ""
	}

	return trimmed[:n]
}

func formattedCard(header model.CardHeader, raw []string) model.Card {
	card := model.Card{Header: header}
	var fence string

	for _, text := range raw {
		next := nextFence(text, fence)
		line := model.Line{Original: text}

		switch {
		case next != fence || strings.TrimSpace(text) == "":
			line.Parts = []model.LinePart{{Visible: true, Content: text}}
		case fence != "":
			indent := len(text) - len(strings.TrimLeft(text, " \t"))
			line.Parts = []model.LinePart{
				{Visible: true, Content: text[:indent]},
				{Content: text[indent:]},
			}
		default:
			// Keep the literal backticks surrounding each inline answer.
			start := 0

			for _, match := range reCode.FindAllStringIndex(text, -1) {
				line.Parts = append(line.Parts,
					model.LinePart{Visible: true, Content: text[start : match[0]+1]},
					model.LinePart{Content: text[match[0]+1 : match[1]-1]},
				)
				start = match[1] - 1
			}

			line.Parts = append(line.Parts,
				model.LinePart{Visible: true, Content: text[start:]},
			)
		}

		card.Lines = append(card.Lines, line)
		fence = next
	}

	return card
}

func makeCardHeader(line string) model.CardHeader {
	level := 0

	for strings.HasPrefix(line, "#") {
		level += 1
		line = line[1:]
	}

	return model.CardHeader{
		Level:   level,
		Content: strings.TrimSpace(line),
	}
}

// `inline code`
var reCode = regexp.MustCompile("`(\\S.*?)`")
