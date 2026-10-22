package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"com.github.yzh44yzh/canon_d/internal"
	"github.com/charmbracelet/x/ansi"
)

var (
	activeStyle  = lipgloss.NewStyle().Bold(true).Reverse(true)
	correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	wrongStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Underline(true)
)

type model struct {
	name      string
	session   *internal.Session
	input     textinput.Model
	viewport  viewport.Model
	width     int
	height    int
	selection int
}

func Run(name string, session *internal.Session) error {
	_, err := tea.NewProgram(newModel(name, session)).Run()
	return err
}

func newModel(name string, session *internal.Session) *model {
	input := textinput.New()
	input.SetVirtualCursor(true)
	input.Placeholder = "Type the missing text"
	m := &model{name: name, session: session, input: input, viewport: viewport.New()}
	m.resize(80, 24)
	return m
}

func (m *model) Init() tea.Cmd { return m.input.Focus() }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "pgup":
			m.viewport.PageUp()
			return m, nil
		case "pgdown":
			m.viewport.PageDown()
			return m, nil
		}
		if m.width < 20 || m.height < 10 {
			return m, nil
		}
		if m.session.Done() {
			return m.updateSummary(msg)
		}
		if msg.String() == "enter" {
			if m.session.AwaitingContinue() || m.session.AnswerCount() == 0 {
				before := m.session.CardIndex()
				m.session.Continue()
				if before != m.session.CardIndex() {
					m.viewport.GotoTop()
				}
				m.input.Reset()
			} else {
				m.session.Submit(m.input.Value())
			}
			m.refreshCard()
			return m, m.input.Focus()
		}
	}
	if !m.session.Done() && !m.session.AwaitingContinue() && m.session.AnswerCount() > 0 {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *model) updateSummary(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.session.MissedCount() == 0 {
		if msg.String() == "enter" {
			return m, tea.Quit
		}
		return m, nil
	}
	switch msg.String() {
	case "up", "down", "tab", "shift+tab":
		m.selection = 1 - m.selection
	case "enter":
		if m.selection == 1 {
			return m, tea.Quit
		}
		m.session.Repeat()
		m.input.Reset()
		m.viewport.GotoTop()
		m.refreshCard()
		return m, m.input.Focus()
	}
	return m, nil
}

func (m *model) resize(width, height int) {
	m.width, m.height = width, height
	m.viewport.SetWidth(max(1, width))
	m.viewport.SetHeight(max(1, height-7))
	m.input.SetWidth(max(1, width-3))
	m.refreshCard()
}

func (m *model) refreshCard() {
	if m.session.Done() {
		return
	}
	content, activeRow := renderCard(m.session, max(1, m.width))
	m.viewport.SetContent(content)
	if activeRow >= 0 {
		m.viewport.EnsureVisible(activeRow, 0, 0)
	}
}

// renderCard wraps display text only; the session retains the original answers.
func renderCard(session *internal.Session, width int) (string, int) {
	var rows []string
	answer, activeRow := 0, -1
	for _, line := range session.Card().Lines {
		var text strings.Builder
		activeLineOffset := -1
		for _, part := range line.Parts {
			if part.Visible {
				text.WriteString(displayText(part.Content))
				continue
			}
			correct, revealed := session.AnswerResult(answer)
			if answer == session.AnswerIndex() {
				activeLineOffset = strings.Count(ansi.Hardwrap(text.String()+"[", width, true), "\n")
			}
			if revealed {
				style := wrongStyle
				if correct {
					style = correctStyle
				}
				text.WriteString(style.Render(displayText(part.Content)))
			} else {
				blank := fmt.Sprintf("[%d: ???]", answer+1)
				if answer == session.AnswerIndex() {
					blank = activeStyle.Render(blank)
				}
				text.WriteString(blank)
			}
			answer++
		}
		if activeLineOffset >= 0 {
			activeRow = len(rows) + activeLineOffset
		}
		rows = append(rows, strings.Split(ansi.Hardwrap(text.String(), width, true), "\n")...)
	}
	return strings.Join(rows, "\n"), activeRow
}

// Hide backticks and expand tabs consistently for terminal display.
func displayText(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '`' {
			return -1
		}
		if r < ' ' && r != '\t' {
			return -1
		}
		if r == '\x7f' {
			return -1
		}
		return r
	}, strings.ReplaceAll(s, "\t", "    "))
}

func (m *model) View() tea.View {
	var content string
	switch {
	case m.width < 20 || m.height < 10:
		content = "Resize terminal to at least 20×10. Ctrl+C to quit."
	case m.session.Done():
		content = fmt.Sprintf("Session complete\n\nCards: %d\nMissed cards: %d\n\n", m.session.Total(), m.session.MissedCount())
		if m.session.MissedCount() > 0 {
			choices := []string{"Repeat missed", "Quit"}
			for i, choice := range choices {
				if i == m.selection {
					choice = "> " + choice
				} else {
					choice = "  " + choice
				}
				content += choice + "\n"
			}
			content += "↑/↓ choose · Enter select · Ctrl+C quit"
		} else {
			content += "Enter to quit · Ctrl+C quit"
		}
	default:
		prompt, status := "Enter to continue", "Reading card"
		if m.session.AnswerCount() > 0 {
			status = fmt.Sprintf("Answer %d/%d", m.session.AnswerIndex()+1, m.session.AnswerCount())
			prompt = m.input.View()
			if m.session.AwaitingContinue() {
				prompt = "Enter to continue"
				status = wrongStyle.Render("Incorrect — expected answer revealed above")
				if m.session.Correct() {
					status = correctStyle.Render("Correct")
				}
			}
		}
		content = fmt.Sprintf("%s\n%s\n\n%s\n\n%s\n%s\n%s",
			ansi.Truncate(fmt.Sprintf("%s · Card %d/%d", displayText(m.name), m.session.CardIndex()+1, m.session.Total()), m.width, "…"),
			ansi.Truncate(displayText(m.session.Card().Header.Show()), m.width, "…"),
			m.viewport.View(), prompt, ansi.Truncate(status, m.width, "…"),
			ansi.Truncate("PgUp/PgDown scroll · Enter submit/continue · Ctrl+C quit", m.width, "…"))
	}
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = ansi.Truncate(line, max(1, m.width), "…")
	}
	content = strings.Join(lines, "\n")
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
