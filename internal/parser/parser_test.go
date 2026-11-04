package parser

import (
	"fmt"
	"strings"
	"testing"

	"com.github.yzh44yzh/canon_d/internal/model"
	"github.com/google/go-cmp/cmp"
)

func Test_makeCardHeader(t *testing.T) {
	sets := []struct {
		in    string
		out   string
		level int
	}{
		{in: "# aaa ", out: "aaa", level: 1},
		{in: "#    aaa bb ", out: "aaa bb", level: 1},
		{in: "##  bbbb ", out: "bbbb", level: 2},
		{in: "### ccc", out: "ccc", level: 3},
		{in: "##### Some Header", out: "Some Header", level: 5},
	}

	for i, set := range sets {
		name := fmt.Sprintf("set_%d", i)
		t.Run(name, func(t *testing.T) {
			header := makeCardHeader(set.in)
			if header.Content != set.out {
				t.Errorf("in: '%v', expect out: '%v', real out: '%v'",
					set.in, set.out, header.Content)
			}

			if header.Level != set.level {
				t.Errorf("in: '%v', expect level: '%v', real level: '%v'",
					set.in, set.level, header.Level)
			}
		})
	}
}

func TestCardHeaderShow(t *testing.T) {
	h1 := model.CardHeader{
		Level:   1,
		Content: "aaa",
	}

	h2 := model.CardHeader{
		Level:   2,
		Content: "bbb",
		Parent:  &h1,
	}

	h3 := model.CardHeader{
		Level:   3,
		Content: "ccc",
		Parent:  &h2,
	}

	s1 := h1.Show()
	if s1 != "aaa" {
		t.Errorf("invalid result for h1.Show '%v'", s1)
	}

	s2 := h2.Show()
	if s2 != "aaa / bbb" {
		t.Errorf("invalid result for h2.Show '%v'", s2)
	}

	s3 := h3.Show()
	if s3 != "aaa / bbb / ccc" {
		t.Errorf("invalid result for h3.Show '%v'", s3)
	}
}

func TestPreserveFormatting(t *testing.T) {
	body := "\n  Before `hidden` and `ещё` **visible**.\n\n" +
		"```go\n\tfunc main() {\n\n    # code, not a heading\n\t}\n```\n  After.\n"

	deck := MakeDeck("test", "# Parent\n## Child\n"+body+"\n## Reading\nJust read.\n")
	if len(deck.Cards) != 2 {
		t.Fatalf("want 2 cards, got %d", len(deck.Cards))
	}

	card := deck.Cards[0]
	if got := card.Header.Show(); got != "Parent / Child" {
		t.Errorf("unexpected header: %q", got)
	}

	var originals, restored []string
	var answers []string

	for _, line := range card.Lines {
		originals = append(originals, line.Original)
		var full strings.Builder

		for _, part := range line.Parts {
			full.WriteString(part.Content)
			if !part.Visible {
				answers = append(answers, part.Content)
			}
		}

		restored = append(restored, full.String())
	}

	for name, got := range map[string]string{
		"original": strings.Join(originals, "\n"),
		"restored": strings.Join(restored, "\n"),
	} {
		if diff := cmp.Diff(body, got); diff != "" {
			t.Errorf("%s formatting (-want +got):\n%s", name, diff)
		}
	}

	wantAnswers := []string{"hidden", "ещё", "func main() {", "# code, not a heading", "}"}
	if diff := cmp.Diff(wantAnswers, answers); diff != "" {
		t.Errorf("answers (-want +got):\n%s", diff)
	}
}

func TestFormattedFenceBoundaries(t *testing.T) {
	for _, fence := range []string{"````", "~~~"} {
		t.Run(fence, func(t *testing.T) {
			content := "# One\n" + fence + "text\n# hidden\n```\n" +
				fence + "\n# Two\nvisible\n"

			deck := MakeDeck("test", content)
			if len(deck.Cards) != 2 || len(deck.Cards[0].Lines) != 4 {
				t.Fatalf("unexpected deck: %+v", deck)
			}

			if deck.Cards[0].Lines[1].Parts[1].Content != "# hidden" {
				t.Fatal("heading inside fence must be a hidden answer")
			}
		})
	}
}

func TestEmptyCardsAndCRLF(t *testing.T) {
	content := "# Parent\r\n\r\n## Empty\r\n \t\r\n## Card\r\n  `answer`\r\n"

	deck := MakeDeck("test", content)
	if len(deck.Cards) != 1 || deck.Cards[0].Lines[0].Original != "  `answer`" {
		t.Fatalf("empty cards or CRLF not handled: %+v", deck)
	}

	for _, input := range []string{"", "# Empty\n\n", "no header\n"} {
		if got := MakeDeck("test", input); len(got.Cards) != 0 {
			t.Errorf("want no cards for %q, got %+v", input, got)
		}
	}
}

func TestSourceWalkBehavior(t *testing.T) {
	tests := []struct {
		name, content string
		headers       []string
		bodies        []string
		answers       [][]string
		findings      []Finding
	}{
		{
			name: "preamble fence and unrestricted heading levels",
			content: "```text\n# ignored\n```\n#Root\n#######Child\n" +
				"  # literal\n`ok`\n\n",
			headers: []string{"Root / Child"},
			bodies:  []string{"  # literal\n`ok`\n\n"},
			answers: [][]string{{"ok"}},
			findings: []Finding{
				{
					Line: 1, Warning: true,
					Message: "content before the first heading is ignored",
				},
				{
					Line: 2, Warning: true,
					Message: "content before the first heading is ignored",
				},
				{
					Line: 3, Warning: true,
					Message: "content before the first heading is ignored",
				},
			},
		},
		{
			name:    "unclosed preamble fence hides headings",
			content: "~~~\n# apparent\n",
			findings: []Finding{
				{
					Line: 1, Warning: true,
					Message: "content before the first heading is ignored",
				},
				{Line: 1, Message: "unclosed fenced block"},
				{Line: 1, Message: "empty deck: no usable cards"},
				{
					Line: 2, Warning: true,
					Message: "content before the first heading is ignored",
				},
			},
		},
		{
			name: "equal-line heading findings retain order", content: "#",
			findings: []Finding{
				{Line: 1, Message: "empty heading title"},
				{Line: 1, Message: "empty leaf card"},
				{Line: 1, Message: "empty deck: no usable cards"},
			},
		},
		{
			name: "equal-line preamble findings retain order", content: "```",
			findings: []Finding{
				{
					Line: 1, Warning: true,
					Message: "content before the first heading is ignored",
				},
				{Line: 1, Message: "unclosed fenced block"},
				{Line: 1, Message: "empty fenced block"},
				{Line: 1, Message: "empty deck: no usable cards"},
			},
		},
		{
			name: "LF trailing blank lines", content: "# Read\ntext\n\n",
			headers: []string{"Read"}, bodies: []string{"text\n\n"},
			answers: [][]string{nil},
		},
		{
			name: "CRLF trailing blank lines", content: "# Read\r\ntext\r\n\r\n",
			headers: []string{"Read"}, bodies: []string{"text\n\n"},
			answers: [][]string{nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deck := MakeDeck("test", tt.content)
			var headers, bodies []string
			var answers [][]string

			for _, card := range deck.Cards {
				headers = append(headers, card.Header.Show())
				var originals []string
				var hidden []string

				for _, line := range card.Lines {
					originals = append(originals, line.Original)
					var restored strings.Builder

					for _, part := range line.Parts {
						restored.WriteString(part.Content)

						if !part.Visible {
							hidden = append(hidden, part.Content)
						}
					}

					if restored.String() != line.Original {
						t.Fatalf("source changed: %q", line.Original)
					}
				}

				bodies = append(bodies, strings.Join(originals, "\n"))
				answers = append(answers, hidden)
			}

			if diff := cmp.Diff(tt.headers, headers); diff != "" {
				t.Errorf("headers (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.bodies, bodies); diff != "" {
				t.Errorf("bodies (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.answers, answers); diff != "" {
				t.Errorf("answers (-want +got):\n%s", diff)
			}

			want := Validation{Cards: len(tt.headers), Findings: tt.findings}

			for _, hidden := range tt.answers {
				if len(hidden) == 0 {
					want.ReadingCards++
				}
			}

			for _, finding := range tt.findings {
				if finding.Warning {
					want.Warnings++
				} else {
					want.Errors++
				}
			}

			if diff := cmp.Diff(want, Validate(tt.content)); diff != "" {
				t.Errorf("validation (-want +got):\n%s", diff)
			}
		})
	}
}
