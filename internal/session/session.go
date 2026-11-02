package session

import (
	"errors"
	"math/rand/v2"
	"strings"

	"com.github.yzh44yzh/canon_d/internal/model"
)

// Session owns a TUI practice round independently of terminal rendering.
type Session struct {
	cards            []practiceCard
	cardIndex        int
	answers          []string
	answerIndex      int
	awaitingContinue bool
	correct          bool
	results          []bool
	cardMissed       bool
	missed           []practiceCard
	shown            bool
	statistics       map[model.CardIdentity]model.CardStatistics
}

type practiceCard struct {
	card model.Card
	deck string
}

func NewSession(
	deck model.Deck, limit int,
	history map[model.CardIdentity]model.CardStatistics,
) (*Session, error) {
	return NewMultiSession([]model.Deck{deck}, limit, history)
}

// NewMultiSession orders selected cards by saved priorities before one global limit.
func NewMultiSession(
	decks []model.Deck, limit int,
	history map[model.CardIdentity]model.CardStatistics,
) (*Session, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be greater than zero")
	}

	var cards []practiceCard

	for _, deck := range decks {
		for _, card := range deck.Cards {
			cards = append(cards, practiceCard{card: card, deck: deck.Header})
		}
	}

	if len(cards) == 0 {
		return nil, errors.New("selected decks contain no cards")
	}

	weightedShuffle(cards, history, rand.ExpFloat64)

	s := &Session{
		cards:      cards[:min(limit, len(cards))],
		statistics: make(map[model.CardIdentity]model.CardStatistics),
	}
	s.startCard()

	return s, nil
}

func shuffleCards(cards []practiceCard) {
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
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

func (s *Session) Card() model.Card {
	if s.Done() {
		return model.Card{}
	}

	return s.cards[s.cardIndex].card
}

// DeckName identifies the current card's source, including in retry rounds.
func (s *Session) DeckName() string {
	if s.Done() {
		return ""
	}

	return s.cards[s.cardIndex].deck
}

func (s *Session) startCard() {
	s.shown = false
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
	s.recordAnswer()
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
		s.missed = append(s.missed, s.cards[s.cardIndex])
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
