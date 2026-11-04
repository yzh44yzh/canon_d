package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"com.github.yzh44yzh/canon_d/internal/decks"
	domain "com.github.yzh44yzh/canon_d/internal/model"
	"com.github.yzh44yzh/canon_d/internal/session"
)

type deckPicker struct {
	entries  []decks.Entry
	selected []bool
	cursor   int
	offset   int
	limit    int
	history  map[domain.CardIdentity]domain.CardStatistics
}

func RunSelection(
	entries []decks.Entry, limit int,
	history map[domain.CardIdentity]domain.CardStatistics,
	save func(*session.Session) error,
) error {
	return runModel(newSelectionModel(entries, limit, history, DefaultColorScheme), save)
}

func newSelectionModel(
	entries []decks.Entry, limit int,
	history map[domain.CardIdentity]domain.CardStatistics,
	scheme ColorScheme,
) *model {
	m := newModel(nil, scheme)
	m.picker = &deckPicker{
		entries: entries, selected: make([]bool, len(entries)), limit: limit,
		history: history,
	}

	return m
}

func (m *model) pickerRows() int { return max(1, m.height-5) }

func (p *deckPicker) ensureVisible(rows int) {
	p.offset = min(p.offset, max(0, len(p.entries)-rows))
	if p.cursor < p.offset {
		p.offset = p.cursor
	}

	if p.cursor >= p.offset+rows {
		p.offset = p.cursor - rows + 1
	}
}

func (p *deckPicker) selectedDecks() []domain.Deck {
	var selected []domain.Deck

	for i, entry := range p.entries {
		if p.selected[i] {
			selected = append(selected, entry.Deck)
		}
	}

	return selected
}

func (m *model) updatePicker(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	p := m.picker
	if len(p.entries) == 0 {
		return m, nil
	}

	switch msg.String() {
	case "up":
		p.cursor = max(0, p.cursor-1)
	case "down":
		p.cursor = min(len(p.entries)-1, p.cursor+1)
	case "pgup":
		p.cursor = max(0, p.cursor-m.pickerRows())
	case "pgdown":
		p.cursor = min(len(p.entries)-1, p.cursor+m.pickerRows())
	case "space":
		if p.entries[p.cursor].Err == nil {
			p.selected[p.cursor] = !p.selected[p.cursor]
		}
	case "enter":
		selected := p.selectedDecks()
		if len(selected) == 0 {
			return m, nil
		}

		s, err := session.NewMultiSession(selected, p.limit, p.history)
		if err != nil {
			return m, nil
		}

		m.practice.start(s)
		m.picker = nil
		m.refreshCard()

		return m, m.input.Focus()
	}

	p.ensureVisible(m.pickerRows())

	return m, nil
}

func (m *model) pickerView() string {
	p := m.picker
	var rows []string
	end := min(len(p.entries), p.offset+m.pickerRows())

	for i := p.offset; i < end; i++ {
		entry := p.entries[i]
		mark, prefix := "[ ]", "  "
		if p.selected[i] {
			mark = "[x]"
		}

		if i == p.cursor {
			prefix = "> "
		}

		text := fmt.Sprintf("%s%s %s · %d cards",
			prefix, mark, displayText(entry.Deck.Header), len(entry.Deck.Cards))
		if entry.Err != nil {
			text = fmt.Sprintf("%s[-] %s · %s",
				prefix, displayText(entry.Deck.Header),
				displayText(entry.Reason()))
			text = m.styles.muted.Render(text)
		}

		if i == p.cursor {
			text = m.styles.active.Render(text)
		}

		rows = append(rows, text)
	}

	return strings.Join([]string{
		"Select Decks",
		fmt.Sprintf("Selected decks: %d", len(p.selectedDecks())),
		m.panel(strings.Join(rows, "\n")),
		m.styles.muted.Render(
			"↑/↓ move · Space toggle · Enter learn · PgUp/PgDown · Ctrl+C quit",
		),
	}, "\n")
}
