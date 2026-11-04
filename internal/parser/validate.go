package parser

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"com.github.yzh44yzh/canon_d/internal/model"
)

// Finding identifies a source problem. RelatedLine locates the first duplicate.
type Finding struct {
	Line        int
	RelatedLine int
	Warning     bool
	Message     string
}

// Validation reports authoring problems without changing permissive deck parsing.
type Validation struct {
	Cards        int
	ReadingCards int
	Warnings     int
	Errors       int
	Findings     []Finding
}

// Validate checks the stricter authoring rules while using practice card counts.
func Validate(content string) Validation {
	source := scanSource(content)
	deck := source.deck("")
	result := Validation{Cards: len(deck.Cards)}

	for _, card := range deck.Cards {
		if isReadingCard(card) {
			result.ReadingCards++
		}
	}

	result.checkLines(source.preamble)

	for _, heading := range source.headings {
		if heading.header.Content == "" {
			result.add(heading.line, "empty heading title")
		}

		if strings.Contains(heading.header.Content, "/") {
			result.add(heading.line, "slash is not allowed in heading titles")
		}

		result.checkLines(heading.lines)
	}

	result.checkHeadings(source.headings)

	if result.Cards == 0 {
		result.add(1, "empty deck: no usable cards")
	}

	sort.SliceStable(result.Findings, func(i, j int) bool {
		return result.Findings[i].Line < result.Findings[j].Line
	})

	for _, finding := range result.Findings {
		if finding.Warning {
			result.Warnings++
		} else {
			result.Errors++
		}
	}

	return result
}

// checkLines applies authoring rules to the source walk's fence classifications.
// A fence cannot cross a recognized heading, so each body can be checked separately.
func (v *Validation) checkLines(lines []sourceLine) {
	var fenceLine int
	var fenceBody bool
	var unclosed bool

	for _, line := range lines {
		if line.preamble && !line.blank {
			v.Findings = append(v.Findings, Finding{
				Line: line.number, Warning: true,
				Message: "content before the first heading is ignored",
			})
		}

		switch {
		case line.fenceBefore == "" && line.fenceAfter != "":
			fenceLine, fenceBody = line.number, false
		case line.fenceBefore != "" && line.fenceAfter == "":
			if !fenceBody {
				v.add(fenceLine, "empty fenced block")
			}
		case line.fenceBefore != "":
			fenceBody = fenceBody || !line.blank
		default:
			if message := inlineProblem(line.text); message != "" {
				v.add(line.number, message)
			}
		}

		unclosed = line.fenceAfter != ""
	}

	if unclosed {
		v.add(fenceLine, "unclosed fenced block")

		if !fenceBody {
			v.add(fenceLine, "empty fenced block")
		}
	}
}

func (v *Validation) add(line int, message string) {
	v.Findings = append(v.Findings, Finding{Line: line, Message: message})
}

func (v *Validation) checkHeadings(headings []*sourceHeading) {
	seen := make(map[string]int)

	for _, heading := range headings {
		if !heading.body {
			if !heading.children {
				v.add(heading.line, "empty leaf card")
			}

			continue
		}

		path := heading.header.Show()

		if first, ok := seen[path]; ok {
			v.Findings = append(v.Findings, Finding{
				Line: heading.line, RelatedLine: first,
				Message: "duplicate card path: " + path,
			})
		} else {
			seen[path] = heading.line
		}
	}
}

func isReadingCard(card model.Card) bool {
	for _, line := range card.Lines {
		for _, part := range line.Parts {
			if !part.Visible {
				return false
			}
		}
	}

	return true
}

func inlineProblem(line string) string {
	open := -1

	for i := 0; i < len(line); i++ {
		if line[i] != '`' {
			continue
		}

		if i+1 < len(line) && line[i+1] == '`' {
			return "invalid inline answer: multiple-backtick delimiters or empty span"
		}

		if open < 0 {
			open = i
			continue
		}

		first, _ := utf8.DecodeRuneInString(line[open+1 : i])
		if unicode.IsSpace(first) {
			return "invalid inline answer: answer starts with whitespace"
		}

		open = -1
	}

	if open >= 0 {
		return "invalid inline answer: unmatched backtick"
	}

	return ""
}
