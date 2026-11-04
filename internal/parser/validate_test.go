package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name, content  string
		cards, reading int
		findings       []Finding
	}{
		{
			name:    "sections and reading cards",
			content: "# Go\n## Basics\nread\n## Answers\nUse `Go` and `fmt`.\n",
			cards:   2, reading: 1,
		},
		{
			name:     "empty deck",
			findings: []Finding{{Line: 1, Message: "empty deck: no usable cards"}},
		},
		{
			name:    "empty leaves with skipped levels",
			content: "# Root\n### Empty\n \n## End",
			findings: []Finding{
				{Line: 1, Message: "empty deck: no usable cards"},
				{Line: 2, Message: "empty leaf card"},
				{Line: 4, Message: "empty leaf card"},
			},
		},
		{
			name: "empty title and slash in section", content: "# A/B\n## \nread",
			cards: 1, reading: 1,
			findings: []Finding{
				{Line: 1, Message: "slash is not allowed in heading titles"},
				{Line: 2, Message: "empty heading title"},
			},
		},
		{
			name:    "duplicate trimmed full paths",
			content: "# A\n## B\nread\n## B \nread",
			cards:   2, reading: 2,
			findings: []Finding{
				{Line: 4, RelatedLine: 2, Message: "duplicate card path: A / B"},
			},
		},
		{
			name:    "distinct paths and case",
			content: "# A\n## B\nread\n## b\nread\n# C\n## B\nread",
			cards:   3, reading: 3,
		},
		{
			name:    "repeated sections are not cards",
			content: "# A\n## B\nread\n# A\n## C\nread",
			cards:   2, reading: 2,
		},
		{
			name: "ignored preamble", content: "intro\n\n# Read\ntext",
			cards: 1, reading: 1,
			findings: []Finding{
				{
					Line: 1, Warning: true,
					Message: "content before the first heading is ignored",
				},
			},
		},
		{
			name:     "empty block with other answers",
			content:  "# Card\n`answer`\n```go\n \t\n```\n~~~\nvalue\n~~~",
			cards:    1,
			findings: []Finding{{Line: 3, Message: "empty fenced block"}},
		},
		{
			name: "unclosed block", content: "# Card\n~~~\n# code\n` code",
			cards:    1,
			findings: []Finding{{Line: 2, Message: "unclosed fenced block"}},
		},
		{
			name: "unclosed empty block", content: "# Card\n```\n \t",
			cards: 1, reading: 1,
			findings: []Finding{
				{Line: 2, Message: "unclosed fenced block"},
				{Line: 2, Message: "empty fenced block"},
			},
		},
		{
			name:    "longer fences preserve shorter code markers",
			content: "# Card\n ````go\n```\n# code / `\n~~~~\n ````` ",
			cards:   1,
		},
		{
			name:    "CRLF and optional heading space",
			content: "#Root\r\n##Child\r\n`ok`\r\n",
			cards:   1,
		},
		{
			name:    "source order for delayed findings",
			content: "# Empty\n# Card\n` bad`\n~~~\nvalue",
			cards:   1,
			findings: []Finding{
				{Line: 1, Message: "empty leaf card"},
				{
					Line: 3,
					Message: "invalid inline answer: " +
						"answer starts with whitespace",
				},
				{Line: 4, Message: "unclosed fenced block"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := Validation{
				Cards: tt.cards, ReadingCards: tt.reading, Findings: tt.findings,
			}

			for _, finding := range tt.findings {
				if finding.Warning {
					want.Warnings++
				} else {
					want.Errors++
				}
			}

			if diff := cmp.Diff(want, Validate(tt.content)); diff != "" {
				t.Fatalf("validation mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestValidateInlineAnswers(t *testing.T) {
	tests := []struct {
		name, line string
		valid      bool
	}{
		{name: "prose", line: "reading only", valid: true},
		{name: "multiple answers", line: "`one` and `two`", valid: true},
		{name: "trailing whitespace", line: "`answer \t`", valid: true},
		{name: "unmatched", line: "`answer"},
		{name: "unmatched after answer", line: "`answer` then `"},
		{name: "empty", line: "``"},
		{name: "whitespace only", line: "` \t`"},
		{name: "leading space", line: "` answer`"},
		{name: "leading tab", line: "`\tanswer`"},
		{name: "leading unicode whitespace", line: "`\u2003answer`"},
		{name: "double delimiter", line: "``answer``"},
		{name: "multiline", line: "`first\nsecond`"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := Validate("# Card\n" + tt.line)

			if (report.Errors == 0) != tt.valid {
				t.Fatalf("valid = %v, report = %+v", tt.valid, report)
			}
		})
	}
}
