package parser

import (
	"strings"

	"com.github.yzh44yzh/canon_d/internal/model"
)

// sourceDocument keeps structural recognition shared by practice and validation.
// Bodyless headings remain available for ancestry and authoring diagnostics.
type sourceDocument struct {
	preamble []sourceLine
	headings []*sourceHeading
}

type sourceHeading struct {
	header   model.CardHeader
	line     int
	body     bool
	children bool
	lines    []sourceLine
}

type sourceLine struct {
	text        string
	number      int
	blank       bool
	preamble    bool
	fenceBefore string
	fenceAfter  string
}

// scanSource recognizes headings and fences once, retaining every body line.
func scanSource(content string) sourceDocument {
	var source sourceDocument
	var stack []*sourceHeading
	var current *sourceHeading
	var fence string

	content = strings.ReplaceAll(content, "\r\n", "\n")

	for i, text := range strings.Split(content, "\n") {
		if fence == "" && strings.HasPrefix(text, "#") {
			header := makeCardHeader(text)

			for len(stack) > 0 && stack[len(stack)-1].header.Level >= header.Level {
				stack = stack[:len(stack)-1]
			}

			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.children = true
				header.Parent = &parent.header
			}

			current = &sourceHeading{header: header, line: i + 1}
			source.headings = append(source.headings, current)
			stack = append(stack, current)

			continue
		}

		line := sourceLine{
			text: text, number: i + 1, blank: strings.TrimSpace(text) == "",
			preamble:    current == nil,
			fenceBefore: fence, fenceAfter: nextFence(text, fence),
		}
		fence = line.fenceAfter

		if current == nil {
			source.preamble = append(source.preamble, line)
		} else {
			current.lines = append(current.lines, line)
			current.body = current.body || !line.blank
		}
	}

	return source
}

func (source sourceDocument) deck(header string) model.Deck {
	deck := model.Deck{Header: header}

	for _, heading := range source.headings {
		if heading.body {
			card := formattedCard(heading.header, heading.lines)
			deck.Cards = append(deck.Cards, card)
		}
	}

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
