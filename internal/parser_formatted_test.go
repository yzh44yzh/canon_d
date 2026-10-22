package internal

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPreserveFormatting(t *testing.T) {
	body := "\n  Before `hidden` and `ещё` **visible**.\n\n```go\n\tfunc main() {\n\n    # code, not a heading\n\t}\n```\n  After.\n"
	deck := MakeDeckWithMode("test", "# Parent\n## Child\n"+body+"\n## Reading\nJust read.\n", PreserveFormatting)
	if len(deck.Cards) != 2 {
		t.Fatalf("want 2 cards, got %d", len(deck.Cards))
	}
	card := deck.Cards[0]
	if got := card.Header.Show(); got != "Parent / Child" {
		t.Errorf("unexpected header: %q", got)
	}
	var originals, restored, masked []string
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
		masked = append(masked, line.Show())
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
	wantMask := "\n  Before `???` and `???` **visible**.\n\n```go\n\t???\n\n    ???\n\t???\n```\n  After.\n"
	if diff := cmp.Diff(wantMask, strings.Join(masked, "\n")); diff != "" {
		t.Errorf("masked formatting (-want +got):\n%s", diff)
	}
}

func TestFormattedFenceBoundaries(t *testing.T) {
	for _, fence := range []string{"````", "~~~"} {
		t.Run(fence, func(t *testing.T) {
			content := "# One\n" + fence + "text\n# hidden\n```\n" + fence + "\n# Two\nvisible\n"
			deck := MakeDeckWithMode("test", content, PreserveFormatting)
			if len(deck.Cards) != 2 || len(deck.Cards[0].Lines) != 4 {
				t.Fatalf("unexpected deck: %+v", deck)
			}
			if deck.Cards[0].Lines[1].Parts[1].Content != "# hidden" {
				t.Fatal("heading inside fence must be a hidden answer")
			}
		})
	}
}

func TestParserModes(t *testing.T) {
	content := "# Parent\r\n\r\n## Empty\r\n \t\r\n## Card\r\n  `answer`\r\n"
	if diff := cmp.Diff(MakeDeck("test", content), MakeDeckWithMode("test", content, Compact)); diff != "" {
		t.Errorf("compact mode changed: %s", diff)
	}
	deck := MakeDeckWithMode("test", content, PreserveFormatting)
	if len(deck.Cards) != 1 || deck.Cards[0].Lines[0].Original != "  `answer`" {
		t.Fatalf("empty cards or CRLF not handled: %+v", deck)
	}
	for _, input := range []string{"", "# Empty\n\n", "no header\n"} {
		if got := MakeDeckWithMode("test", input, PreserveFormatting); len(got.Cards) != 0 {
			t.Errorf("want no cards for %q, got %+v", input, got)
		}
	}
}
