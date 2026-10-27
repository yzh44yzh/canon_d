package internal

import (
	"regexp"
	"strings"
)

// MakeDeck parses Markdown cards while preserving their formatting.
func MakeDeck(header, content string) Deck {
	deck := Deck{Header: header}
	var current *CardHeader
	var lines []string
	var fence string

	finish := func() {
		if current == nil || strings.TrimSpace(strings.Join(lines, "\n")) == "" {
			return
		}

		deck.Cards = append(deck.Cards, formattedCard(*current, lines))
	}

	content = strings.ReplaceAll(content, "\r\n", "\n")

	for _, line := range strings.Split(strings.TrimSuffix(content, "\n"), "\n") {
		if fence == "" && strings.HasPrefix(line, "#") {
			finish()
			lines = nil
			h := makeCardHeader(line)

			for current != nil && current.Level >= h.Level {
				current = current.Parent
			}

			h.Parent = current
			current = &h
			continue
		}

		lines = append(lines, line)
		fence = nextFence(line, fence)
	}

	finish()

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

func formattedCard(header CardHeader, raw []string) Card {
	card := Card{Header: header}
	var fence string

	for _, text := range raw {
		next := nextFence(text, fence)
		line := Line{Original: text}

		switch {
		case next != fence || strings.TrimSpace(text) == "":
			line.Parts = []LinePart{{Visible: true, Content: text}}
		case fence != "":
			indent := len(text) - len(strings.TrimLeft(text, " \t"))
			line.Parts = []LinePart{
				{Visible: true, Content: text[:indent]},
				{Content: text[indent:]},
			}
		default:
			// Keep the literal backticks surrounding each inline answer.
			start := 0

			for _, match := range reCode.FindAllStringIndex(text, -1) {
				line.Parts = append(line.Parts,
					LinePart{Visible: true, Content: text[start : match[0]+1]},
					LinePart{Content: text[match[0]+1 : match[1]-1]},
				)
				start = match[1] - 1
			}

			line.Parts = append(line.Parts,
				LinePart{Visible: true, Content: text[start:]},
			)
		}

		card.Lines = append(card.Lines, line)
		fence = next
	}

	return card
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

// `inline code`
var reCode = regexp.MustCompile("`(\\S.*?)`")
