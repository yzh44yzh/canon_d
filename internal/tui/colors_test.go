package tui

import (
	"image/color"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"com.github.yzh44yzh/canon_d/internal"
	"github.com/charmbracelet/x/ansi"
)

func TestColorSchemeReachesUI(t *testing.T) {
	scheme := ColorScheme{
		Border:    lipgloss.Color("3"),
		Active:    lipgloss.Color("5"),
		Muted:     lipgloss.Color("4"),
		Correct:   lipgloss.Color("6"),
		Incorrect: lipgloss.Color("9"),
	}

	deck := internal.MakeDeck("test", "# Card\n`first` then `second`\n")
	session, err := internal.NewSession(deck, 10)
	if err != nil {
		t.Fatal(err)
	}

	m := newModel(deck.Header, session, scheme)
	m.Init()
	wantStyles := scheme.styles()
	check := func(want string) {
		t.Helper()

		if !strings.Contains(m.View().Content, want) {
			t.Errorf("styled fragment %q missing from view", want)
		}
	}

	check(wantStyles.active.Render("[1: ???]"))
	check(wantStyles.muted.Render("Answer 1/2"))
	check(wantStyles.muted.Render("PgUp/PgDown scroll · Enter submit/continue · Ctrl+C quit"))
	if m.styles.panel.GetBorderTopForeground() != scheme.Border {
		t.Error("panel border does not use scheme color")
	}

	input := m.input.Styles()

	for _, state := range []struct {
		name  string
		style lipgloss.Style
		want  color.Color
	}{
		{"focused prompt", input.Focused.Prompt, scheme.Active},
		{"blurred prompt", input.Blurred.Prompt, scheme.Active},
		{"focused placeholder", input.Focused.Placeholder, scheme.Muted},
		{"blurred placeholder", input.Blurred.Placeholder, scheme.Muted},
	} {
		if state.style.GetForeground() != state.want {
			t.Errorf("%s does not use scheme color", state.name)
		}
	}

	if input.Cursor.Color != scheme.Active {
		t.Error("input cursor does not use active color")
	}

	if input.Focused.Text.GetForeground() != lipgloss.NewStyle().GetForeground() {
		t.Error("typed text should inherit terminal color")
	}

	m.input.SetValue("first")
	press(m, tea.KeyEnter)
	check(wantStyles.correct.Render("first"))
	check(wantStyles.correct.Render("Correct"))

	press(m, tea.KeyEnter)
	m.input.SetValue("wrong")
	press(m, tea.KeyEnter)
	check(wantStyles.incorrect.Render("second"))
	check(wantStyles.incorrect.Render("Incorrect — expected answer revealed above"))

	press(m, tea.KeyEnter)
	check(wantStyles.active.Render("> Repeat missed"))

	press(m, tea.KeyDown)
	check(wantStyles.active.Render("> Quit"))
	if !strings.Contains(ansi.Strip(m.View().Content), "  Repeat missed") {
		t.Error("unselected choice should keep its selection marker spacing")
	}
}
