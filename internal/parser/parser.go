package parser

import (
	"regexp"
	"strings"

	"com.github.yzh44yzh/canon_d/internal/model"
)

// MakeDeck parses Markdown cards while preserving their formatting.
func MakeDeck(header, content string) model.Deck {
	return scanSource(content).deck(header)
}

func formattedCard(header model.CardHeader, raw []sourceLine) model.Card {
	card := model.Card{Header: header}

	for _, source := range raw {
		text := source.text
		line := model.Line{Original: text}

		switch {
		case source.fenceBefore != source.fenceAfter || source.blank:
			line.Parts = []model.LinePart{{Visible: true, Content: text}}
		case source.fenceBefore != "":
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
					model.LinePart{
						Visible: true,
						Content: text[start : match[0]+1],
					},
					model.LinePart{Content: text[match[0]+1 : match[1]-1]},
				)
				start = match[1] - 1
			}

			line.Parts = append(line.Parts,
				model.LinePart{Visible: true, Content: text[start:]},
			)
		}

		card.Lines = append(card.Lines, line)
	}

	return card
}

// `inline code`
var reCode = regexp.MustCompile("`(\\S.*?)`")
