package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"com.github.yzh44yzh/canon_d/internal/decks"
	"com.github.yzh44yzh/canon_d/internal/parser"
	"github.com/charmbracelet/x/ansi"
)

func pickerEntry(name string) decks.Entry {
	return decks.Entry{Deck: parser.MakeDeck(name, "# Card\n`answer`")}
}

func TestPickerSelectionToMixedPractice(t *testing.T) {
	entries := []decks.Entry{
		pickerEntry("a.md"),
		{Deck: parser.MakeDeck("empty.md", "# Empty"), Err: fmt.Errorf("no cards")},
		pickerEntry("b.md"),
	}
	m := newSelectionModel(entries, 10, nil, DefaultColorScheme)
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if m.Init() != nil {
		t.Fatal("selection must not focus the answer input")
	}

	view := ansi.Strip(m.View().Content)
	for _, want := range []string{"Select Decks", "Selected decks: 0", "[-] empty.md"} {
		if !strings.Contains(view, want) {
			t.Fatalf("missing %q in selection:\n%s", want, view)
		}
	}

	press(m, tea.KeyEnter)
	if m.session != nil {
		t.Fatal("Enter with no selection must not start practice")
	}

	press(m, tea.KeySpace)
	press(m, tea.KeySpace)
	if m.picker.selected[0] {
		t.Fatal("Space must also deselect")
	}

	press(m, tea.KeySpace)
	press(m, tea.KeyDown)
	press(m, tea.KeySpace)
	if m.picker.selected[1] {
		t.Fatal("unusable decks must remain disabled")
	}

	press(m, tea.KeyDown)
	press(m, tea.KeySpace)
	if !strings.Contains(ansi.Strip(m.View().Content), "Selected decks: 2") {
		t.Fatal("selection count must include both usable decks")
	}

	press(m, tea.KeyEnter)
	if m.picker != nil || m.session.Total() != 2 {
		t.Fatal("Enter should start one combined round")
	}

	for round := 0; round < 2; round++ {
		seen := map[string]bool{}

		for !m.session.Done() {
			name := m.session.DeckName()
			if !strings.Contains(ansi.Strip(m.View().Content), name+" · Card") {
				t.Fatal("header must follow the current card's deck")
			}

			seen[name] = true
			m.input.SetValue("wrong")
			press(m, tea.KeyEnter)
			press(m, tea.KeyEnter)
		}

		if !seen["a.md"] || !seen["b.md"] {
			t.Fatal("both selected decks must survive in full-size and retry rounds")
		}

		press(m, tea.KeyEnter)
	}
}

func TestPickerPagingResizeAndSmallTerminal(t *testing.T) {
	var entries []decks.Entry

	for i := 0; i < 30; i++ {
		entries = append(entries, pickerEntry(fmt.Sprintf("%02d.md", i)))
	}

	m := newSelectionModel(entries, 1, nil, DefaultColorScheme)
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 11})
	press(m, tea.KeyPgDown)
	if m.picker.cursor != 6 || m.picker.offset == 0 {
		t.Fatal("Page Down should move a visible page and scroll the list")
	}

	view := ansi.Strip(m.View().Content)
	if !strings.Contains(view, "> [ ] 06.md") || strings.Contains(view, "00.md") {
		t.Fatalf("unexpected scrolled list:\n%s", view)
	}

	if len(strings.Split(view, "\n")) > 11 ||
		!strings.Contains(view, "Select Decks") || !strings.Contains(view, "Space toggle") {
		t.Fatal("list must fit the terminal with fixed title and hints")
	}

	m.Update(tea.WindowSizeMsg{Width: 19, Height: 10})
	press(m, tea.KeyDown)
	press(m, tea.KeyPgDown)
	press(m, tea.KeySpace)
	if m.picker.cursor != 6 || len(m.picker.selectedDecks()) != 0 {
		t.Fatal("small terminals must suspend selection and navigation")
	}

	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if m.picker.cursor < m.picker.offset ||
		m.picker.cursor >= m.picker.offset+m.pickerRows() {
		t.Fatal("resizing must keep the focused deck visible")
	}

	press(m, tea.KeyPgUp)
	if m.picker.cursor != 0 {
		t.Fatal("Page Up must stop at the first deck")
	}

	m.Update(tea.WindowSizeMsg{Width: 80, Height: 11})
	for i := 0; i < 40; i++ {
		press(m, tea.KeyDown)
	}

	if m.picker.cursor != 29 || !strings.Contains(m.View().Content, "29.md") {
		t.Fatal("Down must stop at the last deck and keep it visible")
	}
}
