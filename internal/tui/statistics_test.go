package tui

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
	"com.github.yzh44yzh/canon_d/internal/decks"
	"com.github.yzh44yzh/canon_d/internal/history"
	domain "com.github.yzh44yzh/canon_d/internal/model"
	"com.github.yzh44yzh/canon_d/internal/parser"
	"com.github.yzh44yzh/canon_d/internal/session"
	"github.com/google/go-cmp/cmp"
)

func connectHistory(t *testing.T, m *model) string {
	t.Helper()

	dir := t.TempDir()
	store, err := history.Open(dir)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { store.Close() })
	m.practice.save = func(s *session.Session) error { return store.Save(s.Statistics()) }

	return filepath.Join(dir, ".canon_d.csv")
}

func checkSavedRow(t *testing.T, path string, row []string) {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	want := [][]string{
		{"deck", "card_header", "shows", "correct_answers", "incorrect_answers"},
		row,
	}
	if diff := cmp.Diff(want, rows); diff != "" {
		t.Fatalf("wrong saved statistics (-want +got):\n%s", diff)
	}
}

func TestCtrlCSavesPartialCardDuringFeedback(t *testing.T) {
	m := testModel(t, "# Card\n`one` and `two`")
	path := connectHistory(t, m)
	m.input.SetValue("one")
	press(m, tea.KeyEnter)
	m.Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	if m.practice.err != nil {
		t.Fatal(m.practice.err)
	}

	checkSavedRow(t, path, []string{"test", "Card", "1", "1", "0"})
}

func TestShowWaitsForVisibleCardAndIgnoresRedraws(t *testing.T) {
	deck := parser.MakeDeck("read.md", "# Read\nReading")
	s, err := session.NewSession(deck, 1, nil)
	if err != nil {
		t.Fatal(err)
	}

	m := newModel(s, DefaultColorScheme)
	m.View()
	m.Update(tea.WindowSizeMsg{Width: 19, Height: 10})
	m.View()
	if len(s.Statistics()) != 0 {
		t.Fatal("cards hidden by an unknown or small terminal must not count")
	}

	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	for i := 0; i < 3; i++ {
		m.View()
		press(m, tea.KeyPgDown)
		m.Update(tea.WindowSizeMsg{Width: 70, Height: 20})
	}

	want := map[domain.CardIdentity]domain.CardStatistics{
		{Deck: "read.md", Header: "Read"}: {Shows: 1},
	}
	if diff := cmp.Diff(want, s.Statistics()); diff != "" {
		t.Fatalf("visible presentation statistics: %s", diff)
	}
}

func TestSelectionHasNoShowsBeforePractice(t *testing.T) {
	m := newSelectionModel([]decks.Entry{pickerEntry("a.md")}, 1, nil, DefaultColorScheme)
	saves := 0
	m.practice.save = func(*session.Session) error {
		saves++
		return nil
	}
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m.Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	if saves != 0 {
		t.Fatal("quitting in selection must not request a checkpoint")
	}

	press(m, tea.KeySpace)
	press(m, tea.KeyEnter)
	want := map[domain.CardIdentity]domain.CardStatistics{
		{Deck: "a.md", Header: "Card"}: {Shows: 1},
	}
	if diff := cmp.Diff(want, m.practice.session.Statistics()); diff != "" {
		t.Fatalf("starting practice must record its first card: %s", diff)
	}
}

func TestSaveFailureStopsPractice(t *testing.T) {
	m := testModel(t, "# Read\nReading")
	want := errors.New("disk unavailable")
	m.practice.save = func(*session.Session) error { return want }
	cmd := press(m, tea.KeyEnter)
	if !errors.Is(m.practice.err, want) || cmd == nil {
		t.Fatal("failed save must retain the error and stop practice")
	}

	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("failed save must quit the TUI")
	}
}
