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

type sourceHeading struct {
	header   model.CardHeader
	line     int
	body     bool
	children bool
}

// Validate checks the stricter authoring rules while using practice card counts.
func Validate(content string) Validation {
	deck := MakeDeck("", content)
	result := Validation{Cards: len(deck.Cards)}

	for _, card := range deck.Cards {
		if isReadingCard(card) {
			result.ReadingCards++
		}
	}

	var headings []*sourceHeading
	var stack []*sourceHeading
	var current *sourceHeading
	var fence string
	var fenceLine int
	var fenceBody bool

	for i, line := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		lineNumber := i + 1

		if fence == "" && strings.HasPrefix(line, "#") {
			header := makeCardHeader(line)

			for len(stack) > 0 && stack[len(stack)-1].header.Level >= header.Level {
				stack = stack[:len(stack)-1]
			}

			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.children = true
				header.Parent = &parent.header
			}

			current = &sourceHeading{header: header, line: lineNumber}
			headings = append(headings, current)
			stack = append(stack, current)

			if header.Content == "" {
				result.add(lineNumber, "empty heading title")
			}

			if strings.Contains(header.Content, "/") {
				result.add(lineNumber, "slash is not allowed in heading titles")
			}

			continue
		}

		if strings.TrimSpace(line) != "" {
			if current == nil {
				result.Findings = append(result.Findings, Finding{
					Line: lineNumber, Warning: true,
					Message: "content before the first heading is ignored",
				})
			} else {
				current.body = true
			}
		}

		next := nextFence(line, fence)

		switch {
		case fence == "" && next != "":
			fenceLine, fenceBody = lineNumber, false
		case fence != "" && next == "":
			if !fenceBody {
				result.add(fenceLine, "empty fenced block")
			}
		case fence != "":
			fenceBody = fenceBody || strings.TrimSpace(line) != ""
		default:
			if message := inlineProblem(line); message != "" {
				result.add(lineNumber, message)
			}
		}

		fence = next
	}

	if fence != "" {
		result.add(fenceLine, "unclosed fenced block")

		if !fenceBody {
			result.add(fenceLine, "empty fenced block")
		}
	}

	result.checkHeadings(headings)

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
