package internal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSessionAnswersAndFeedback(t *testing.T) {
	deck := MakeDeckWithMode("test", "# Card\n`TCP` then `DNS`\n", PreserveFormatting)
	s, err := NewSession(deck, 10)
	if err != nil {
		t.Fatal(err)
	}
	if s.Total() != 1 || s.AnswerCount() != 2 {
		t.Fatal("oversized limit must use available cards")
	}
	s.Continue()
	if s.AnswerIndex() != 0 {
		t.Fatal("cannot advance without submitting an answer")
	}
	s.Submit(" \tTCP  ")
	if !s.Correct() || !s.AwaitingContinue() || s.AnswerIndex() != 0 {
		t.Fatal("correct answer must wait for acknowledgement")
	}
	s.Submit("wrong")
	if !s.Correct() {
		t.Fatal("submitting during feedback must have no effect")
	}
	s.Continue()
	if correct, revealed := s.AnswerResult(0); !correct || !revealed {
		t.Fatal("previous answer should stay revealed")
	}
	if _, revealed := s.AnswerResult(1); revealed {
		t.Fatal("next answer must remain hidden")
	}
	s.Submit("dns")
	if s.Correct() {
		t.Fatal("matching must be case-sensitive")
	}
	s.Continue()
	if !s.Done() || s.MissedCount() != 1 {
		t.Fatal("one failed answer must mark the card missed")
	}
	s.Repeat()
	if s.Done() || s.MissedCount() != 0 || s.AnswerIndex() != 0 {
		t.Fatal("repeat must start a fresh round")
	}
	if _, revealed := s.AnswerResult(0); revealed {
		t.Fatal("retry must hide previously revealed answers")
	}
	s.Submit("TCP")
	s.Continue()
	s.Submit("DNS")
	s.Continue()
	if !s.Done() || s.MissedCount() != 0 {
		t.Fatal("correct retry should finish without missed cards")
	}
}

func TestSessionReadingAndValidation(t *testing.T) {
	deck := MakeDeckWithMode("test", "# Read\nVisible text\n", PreserveFormatting)
	for _, limit := range []int{0, -1} {
		if _, err := NewSession(deck, limit); err == nil {
			t.Errorf("expected invalid limit error for %d", limit)
		}
	}
	if _, err := NewSession(Deck{}, 10); err == nil {
		t.Fatal("expected empty deck error")
	}
	s, err := NewSession(deck, 1)
	if err != nil {
		t.Fatal(err)
	}
	s.Submit("anything")
	if s.AwaitingContinue() || s.AnswerCount() != 0 {
		t.Fatal("reading cards have no answers")
	}
	s.Continue()
	s.Continue()
	s.Submit("ignored after completion")
	s.Repeat()
	if !s.Done() || s.MissedCount() != 0 {
		t.Fatal("reading card should complete without mistakes")
	}
}

func TestRepeatKeepsOnlyMissedOccurrences(t *testing.T) {
	card := Card{Header: CardHeader{Content: "duplicate"}, Lines: []Line{makeInvisibleLine("answer")}}
	deck := Deck{Cards: []Card{card, card, card}}
	before := append([]Card(nil), deck.Cards...)
	s, err := NewSession(deck, 3)
	if err != nil {
		t.Fatal(err)
	}
	for i, answer := range []string{"wrong", "answer", "wrong"} {
		s.Submit(answer)
		s.Continue()
		if s.CardIndex() != i+1 {
			t.Fatal("unexpected card position")
		}
	}
	if s.MissedCount() != 2 {
		t.Fatal("track each missed occurrence separately")
	}
	s.Repeat()
	if s.Total() != 2 || s.Done() {
		t.Fatal("repeat should include both missed occurrences")
	}
	if diff := cmp.Diff(before, deck.Cards); diff != "" {
		t.Errorf("session mutated source deck: %s", diff)
	}
}

func TestSessionShufflesACopyAndLimits(t *testing.T) {
	deck := MakeDeckWithMode("test", "# A\n`a`\n# B\n`b`\n# C\n`c`\n", PreserveFormatting)
	before := append([]Card(nil), deck.Cards...)
	s, err := NewSession(deck, 2)
	if err != nil {
		t.Fatal(err)
	}
	if s.Total() != 2 {
		t.Fatal("session must honor smaller limits")
	}
	if diff := cmp.Diff(before, deck.Cards); diff != "" {
		t.Errorf("session mutated source deck: %s", diff)
	}
}
