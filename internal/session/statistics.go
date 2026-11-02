package session

import (
	"maps"

	"com.github.yzh44yzh/canon_d/internal/model"
)

// RecordShow records at most one presentation of the current card occurrence.
// The TUI calls it only when the terminal can display the practice screen.
func (s *Session) RecordShow() {
	if s.Done() || s.shown {
		return
	}

	s.shown = true
	id := s.cardIdentity()
	totals := s.statistics[id]
	totals.Shows++
	s.statistics[id] = totals
}

// Statistics returns cumulative results for this run, including unfinished cards.
func (s *Session) Statistics() map[model.CardIdentity]model.CardStatistics {
	return maps.Clone(s.statistics)
}

func (s *Session) cardIdentity() model.CardIdentity {
	return model.CardIdentity{Deck: s.DeckName(), Header: s.Card().Header.Show()}
}

func (s *Session) recordAnswer() {
	id := s.cardIdentity()
	totals := s.statistics[id]
	if s.correct {
		totals.CorrectAnswers++
	} else {
		totals.IncorrectAnswers++
	}

	s.statistics[id] = totals
}
