package internal

import (
	"fmt"
	"strings"
	"testing"

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
				t.Errorf("in: '%v', expect out: '%v', real out: '%v'", set.in, set.out, header.Content)
			}
			if header.Level != set.level {
				t.Errorf("in: '%v', expect level: '%v', real level: '%v'", set.in, set.level, header.Level)
			}
		})
	}
}

func TestCardHeaderShow(t *testing.T) {
	h1 := CardHeader{
		Level:   1,
		Content: "aaa",
	}
	h2 := CardHeader{
		Level:   2,
		Content: "bbb",
		Parent:  &h1,
	}
	h3 := CardHeader{
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
	body := "\n  Before `hidden` and `ещё` **visible**.\n\n```go\n\tfunc main() {\n\n    # code, not a heading\n\t}\n```\n  After.\n"
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
	for name, got := range map[string]string{"original": strings.Join(originals, "\n"), "restored": strings.Join(restored, "\n")} {
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
			content := "# One\n" + fence + "text\n# hidden\n```\n" + fence + "\n# Two\nvisible\n"
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
