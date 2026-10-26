package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"com.github.yzh44yzh/canon_d/internal"
	"github.com/charmbracelet/x/ansi"
)

func testModel(t *testing.T, content string) *model {
	t.Helper()
	deck := internal.MakeDeck("test", content)
	session, err := internal.NewSession(deck, 10)
	if err != nil {
		t.Fatal(err)
	}
	m := newModel(deck.Header, session, DefaultColorScheme)
	m.Init()
	return m
}

func press(m *model, code rune) tea.Cmd {
	_, cmd := m.Update(tea.KeyPressMsg{Code: code})
	return cmd
}

func TestWholeCardMasksAndRevealsSequentially(t *testing.T) {
	m := testModel(t, "# Card\nBefore `secret-one` then `secret-two`.\n\nTail visible.\n")
	initial := ansi.Strip(m.View().Content)
	for _, want := range []string{"Before [1: ???] then [2: ???].", "Tail visible."} {
		if !strings.Contains(initial, want) {
			t.Errorf("whole card missing %q:\n%s", want, initial)
		}
	}
	if strings.Contains(initial, "`") {
		t.Fatal("TUI should hide backticks")
	}
	if strings.Contains(initial, "secret-") {
		t.Fatal("initial view leaked an answer")
	}
	m.input.SetValue("wrong")
	press(m, tea.KeyEnter)
	feedback := ansi.Strip(m.View().Content)
	if !strings.Contains(feedback, "secret-one") || strings.Contains(feedback, "secret-two") || !strings.Contains(feedback, "Incorrect") {
		t.Fatalf("unexpected feedback:\n%s", feedback)
	}
	press(m, 'q')
	if m.input.Value() != "wrong" || !m.session.AwaitingContinue() {
		t.Fatal("typing during feedback must not submit or advance")
	}
	press(m, tea.KeyEnter)
	if m.session.AnswerIndex() != 1 || m.input.Value() != "" {
		t.Fatal("acknowledgement should advance and clear input")
	}
	m.input.SetValue("secret-two")
	press(m, tea.KeyEnter)
	if !strings.Contains(ansi.Strip(m.View().Content), "Correct") {
		t.Fatal("correct feedback missing")
	}
	press(m, tea.KeyEnter)
	if !m.session.Done() || !strings.Contains(m.View().Content, "Missed cards: 1") {
		t.Fatal("missing session summary")
	}
	press(m, tea.KeyEnter)
	if m.session.Done() || strings.Contains(m.View().Content, "secret-one") {
		t.Fatal("repeat should hide all answers again")
	}
}

func TestTypingAndScrollControls(t *testing.T) {
	m := testModel(t, "# Card\n`q`\n"+strings.Repeat("Visible line\n", 40))
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	if m.input.Value() != "q" {
		t.Fatal("q must be ordinary answer input")
	}
	if cmd != nil {
		if _, quit := cmd().(tea.QuitMsg); quit {
			t.Fatal("typing q must not quit")
		}
	}
	press(m, tea.KeyPgDown)
	if m.viewport.YOffset() == 0 || m.input.Value() != "q" {
		t.Fatal("Page Down should scroll without editing the answer")
	}
	press(m, tea.KeyPgUp)
	if m.viewport.YOffset() != 0 {
		t.Fatal("Page Up should return to the top")
	}
	_, quit := m.Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	if quit == nil {
		t.Fatal("Ctrl+C must quit")
	}
	if _, ok := quit().(tea.QuitMsg); !ok {
		t.Fatal("expected quit command")
	}
}

func TestActiveAnswerScrollsAndResizes(t *testing.T) {
	m := testModel(t, "# Card\n`first`\n"+strings.Repeat("Visible line\n", 30)+strings.Repeat("x", 90)+" `last`\n")
	m.input.SetValue("first")
	press(m, tea.KeyEnter)
	press(m, tea.KeyEnter)
	if m.viewport.YOffset() == 0 || !strings.Contains(ansi.Strip(m.viewport.View()), "[2: ???]") {
		t.Fatalf("active answer should scroll into view:\n%s", m.viewport.View())
	}
	m.Update(tea.WindowSizeMsg{Width: 25, Height: 12})
	if !strings.Contains(strings.Join(strings.Fields(ansi.Strip(m.viewport.View())), ""), "[2:???]") {
		t.Fatalf("wrapped active answer lost after resize:\n%s", m.viewport.View())
	}
	m.input.SetValue("last")
	press(m, tea.KeyEnter)
	if !strings.Contains(strings.Join(strings.Fields(ansi.Strip(m.viewport.View())), ""), "last") {
		t.Fatal("revealed answer should remain visible")
	}
}

func TestReadingCardsAndSummaryChoices(t *testing.T) {
	m := testModel(t, "# Read\nVisible\n")
	press(m, tea.KeyEnter)
	if !m.session.Done() || press(m, tea.KeyEnter) == nil {
		t.Fatal("reading card should advance to summary and allow exit")
	}
	m = testModel(t, "# Test\n`answer`\n")
	press(m, tea.KeyEnter)
	press(m, tea.KeyEnter)
	press(m, tea.KeyDown)
	if press(m, tea.KeyEnter) == nil {
		t.Fatal("Quit choice should exit")
	}
}

func TestViewFitsTerminal(t *testing.T) {
	for _, width := range []int{20, 25, 40, 80} {
		for _, height := range []int{10, 12, 24} {
			m := testModel(t, "# Long title with lots of words\nA visible line with `hidden answer` and some extra long content\n")
			m.Update(tea.WindowSizeMsg{Width: width, Height: height})
			for _, stage := range []string{"input", "feedback", "summary"} {
				switch stage {
				case "feedback":
					m.input.SetValue(strings.Repeat("界", 100))
					press(m, tea.KeyEnter)
				case "summary":
					press(m, tea.KeyEnter)
				}
				view := m.View()
				if !view.AltScreen {
					t.Fatal("TUI should use alternate screen")
				}
				lines := strings.Split(ansi.Strip(view.Content), "\n")
				if len(lines) > height {
					t.Errorf("%s at %dx%d: %d rows exceed height", stage, width, height, len(lines))
				}
				for _, line := range lines {
					if got := ansi.StringWidth(line); got > width {
						t.Errorf("%s: line width %d exceeds %d: %q", stage, got, width, line)
					}
				}
				panels := 2
				if stage == "summary" {
					panels = 1
				}
				for _, border := range []string{"╭" + strings.Repeat("─", width-2) + "╮", "╰" + strings.Repeat("─", width-2) + "╯"} {
					if strings.Count(strings.Join(lines, "\n"), border) != panels {
						t.Errorf("%s at %dx%d: expected %d complete borders %q:\n%s", stage, width, height, panels, border, view.Content)
					}
				}
				if !strings.Contains(lines[len(lines)-1], "quit") && width >= 80 {
					t.Errorf("%s: keyboard hints missing below panels", stage)
				}
			}
		}
	}
}
