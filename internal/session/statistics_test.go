package session

import (
	"testing"

	"com.github.yzh44yzh/canon_d/internal/model"
	"com.github.yzh44yzh/canon_d/internal/parser"
	"github.com/google/go-cmp/cmp"
)

func TestStatisticsTrackAnswersShowsAndRetries(t *testing.T) {
	deck := parser.MakeDeck("a.md", "# Parent\n## Card\n`one` and `two`")
	s, err := NewSession(deck, 10)
	if err != nil {
		t.Fatal(err)
	}

	s.RecordShow()
	s.RecordShow()
	s.Submit("one")
	s.Submit("ignored during feedback")
	s.Continue()
	s.Submit("wrong")
	id := model.CardIdentity{Deck: "a.md", Header: "Parent / Card"}
	want := map[model.CardIdentity]model.CardStatistics{
		id: {Shows: 1, CorrectAnswers: 1, IncorrectAnswers: 1},
	}
	if diff := cmp.Diff(want, s.Statistics()); diff != "" {
		t.Fatalf("unfinished card totals: %s", diff)
	}

	s.Continue()
	s.RecordShow()
	s.Repeat()
	s.RecordShow()
	s.Submit("one")
	s.Continue()
	s.Submit("two")
	s.Continue()
	want[id] = model.CardStatistics{Shows: 2, CorrectAnswers: 3, IncorrectAnswers: 1}
	if diff := cmp.Diff(want, s.Statistics()); diff != "" {
		t.Fatalf("retry totals: %s", diff)
	}

	copy := s.Statistics()
	delete(copy, id)
	if len(s.Statistics()) != 1 {
		t.Fatal("statistics snapshots must not expose mutable session state")
	}
}

func TestDuplicateReadingCardsShareStatistics(t *testing.T) {
	deck := parser.MakeDeck("read.md", "# Same\nReading\n# Same\nMore reading")
	s, err := NewSession(deck, 10)
	if err != nil {
		t.Fatal(err)
	}

	for !s.Done() {
		s.RecordShow()
		s.Submit("ignored on reading cards")
		s.Continue()
	}

	want := map[model.CardIdentity]model.CardStatistics{
		{Deck: "read.md", Header: "Same"}: {Shows: 2},
	}
	if diff := cmp.Diff(want, s.Statistics()); diff != "" {
		t.Fatalf("duplicate reading totals: %s", diff)
	}
}
