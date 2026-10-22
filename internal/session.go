package internal

import (
	"errors"
	"math/rand/v2"
	"strings"
)

// Session owns a TUI practice round independently of terminal rendering.
type Session struct {
	cards            []Card
	cardIndex        int
	answers          []string
	answerIndex      int
	awaitingContinue bool
	correct          bool
	results          []bool
	cardMissed       bool
	missed           []Card
}

func NewSession(deck Deck, limit int) (*Session, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be greater than zero")
	}
	if len(deck.Cards) == 0 {
		return nil, errors.New("deck contains no cards")
	}
	cards := append([]Card(nil), deck.Cards...)
	shuffleCards(cards)
	s := &Session{cards: cards[:min(limit, len(cards))]}
	s.startCard()
	return s, nil
}

func shuffleCards(cards []Card) {
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
}

func (s *Session) Done() bool             { return s.cardIndex == len(s.cards) }
func (s *Session) CardIndex() int         { return s.cardIndex }
func (s *Session) Total() int             { return len(s.cards) }
func (s *Session) AnswerIndex() int       { return s.answerIndex }
func (s *Session) AnswerCount() int       { return len(s.answers) }
func (s *Session) AwaitingContinue() bool { return s.awaitingContinue }
func (s *Session) Correct() bool          { return s.correct }
func (s *Session) MissedCount() int       { return len(s.missed) }

func (s *Session) AnswerResult(index int) (correct, revealed bool) {
	if index < 0 || index >= len(s.results) {
		return false, false
	}
	return s.results[index], true
}

func (s *Session) Card() Card {
	if s.Done() {
		return Card{}
	}
	return s.cards[s.cardIndex]
}

func (s *Session) startCard() {
	s.answers = nil
	s.results = nil
	s.answerIndex = 0
	s.awaitingContinue = false
	s.correct = false
	s.cardMissed = false
	for _, line := range s.Card().Lines {
		for _, part := range line.Parts {
			if !part.Visible {
				s.answers = append(s.answers, part.Content)
			}
		}
	}
}

func (s *Session) Submit(answer string) {
	if s.Done() || s.awaitingContinue || len(s.answers) == 0 {
		return
	}
	s.correct = strings.TrimSpace(answer) == strings.TrimSpace(s.answers[s.answerIndex])
	s.results = append(s.results, s.correct)
	s.cardMissed = s.cardMissed || !s.correct
	s.awaitingContinue = true
}

// Continue acknowledges feedback or advances a reading-only card.
func (s *Session) Continue() {
	if s.Done() || (!s.awaitingContinue && len(s.answers) != 0) {
		return
	}
	s.answerIndex++
	s.awaitingContinue = false
	if s.answerIndex < len(s.answers) {
		return
	}
	if s.cardMissed {
		s.missed = append(s.missed, s.Card())
	}
	s.cardIndex++
	s.startCard()
}

// Repeat starts a newly shuffled round containing each missed card once.
func (s *Session) Repeat() {
	if !s.Done() || len(s.missed) == 0 {
		return
	}
	s.cards = s.missed
	s.missed = nil
	shuffleCards(s.cards)
	s.cardIndex = 0
	s.startCard()
}
