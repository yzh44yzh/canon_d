package session

import (
	"fmt"
	"maps"
	"math"
	"math/rand/v2"
	"testing"

	"com.github.yzh44yzh/canon_d/internal/model"
	"github.com/google/go-cmp/cmp"
)

func selectionCard(header string) practiceCard {
	return practiceCard{deck: "test.md", card: model.Card{
		Header: model.CardHeader{Content: header},
	}}
}

func TestInitialWeightsUseEligibleSavedAnswerTotals(t *testing.T) {
	cards := []practiceCard{
		selectionCard("New"), selectionCard("Shown"), selectionCard("Hard"),
		selectionCard("Easy"), selectionCard("No correct"), selectionCard("Hard"),
	}
	history := map[model.CardIdentity]model.CardStatistics{
		{Deck: "test.md", Header: "Shown"}:      {Shows: 100},
		{Deck: "test.md", Header: "Hard"}:       {CorrectAnswers: 1, IncorrectAnswers: 3},
		{Deck: "test.md", Header: "Easy"}:       {CorrectAnswers: 9},
		{Deck: "test.md", Header: "No correct"}: {IncorrectAnswers: 4},
		{Deck: "other.md", Header: "Excluded"}:  {IncorrectAnswers: math.MaxUint64},
	}
	before := maps.Clone(history)
	got := initialWeights(cards, history)
	want := []float64{10, 10, 2, 0.1, 5, 2}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("wrong weights (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(before, history); diff != "" {
		t.Fatalf("selection mutated saved history: %s", diff)
	}
}

func TestAllNewCardsHaveEqualWeights(t *testing.T) {
	cards := []practiceCard{selectionCard("Unseen"), selectionCard("Reading")}
	history := map[model.CardIdentity]model.CardStatistics{
		{Deck: "test.md", Header: "Reading"}: {Shows: math.MaxUint64},
	}
	if diff := cmp.Diff([]float64{1, 1}, initialWeights(cards, history)); diff != "" {
		t.Fatal(diff)
	}
}

func TestExtremeCountersKeepFinitePositiveWeights(t *testing.T) {
	cards := []practiceCard{
		selectionCard("New"), selectionCard("Errors"), selectionCard("Correct"),
	}
	history := map[model.CardIdentity]model.CardStatistics{
		{Deck: "test.md", Header: "Errors"}:  {IncorrectAnswers: math.MaxUint64},
		{Deck: "test.md", Header: "Correct"}: {CorrectAnswers: math.MaxUint64},
	}
	weights := initialWeights(cards, history)

	for _, weight := range weights {
		if weight <= 0 || math.IsInf(weight, 0) || math.IsNaN(weight) {
			t.Fatalf("invalid weight: %v", weight)
		}
	}

	if weights[0] != 2*weights[1] || weights[1] <= weights[2] {
		t.Fatal("maximum counters must preserve new and old priority ordering")
	}
}

func TestWeightedSelectionProbabilities(t *testing.T) {
	tests := []struct {
		name    string
		a, b    model.CardStatistics
		chanceA float64
	}{
		{name: "all new", chanceA: 0.5},
		{
			name: "new versus old", b: model.CardStatistics{CorrectAnswers: 1},
			chanceA: 2.0 / 3.0,
		},
		{
			name: "old ratio",
			a:    model.CardStatistics{CorrectAnswers: 1, IncorrectAnswers: 3},
			b:    model.CardStatistics{CorrectAnswers: 1}, chanceA: 0.8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			history := map[model.CardIdentity]model.CardStatistics{
				{Deck: "test.md", Header: "A"}: tt.a,
				{Deck: "test.md", Header: "B"}: tt.b,
			}
			random := rand.New(rand.NewPCG(42, 7))
			firstA := 0
			const trials = 10000

			for i := 0; i < trials; i++ {
				cards := []practiceCard{selectionCard("A"), selectionCard("B")}
				weightedShuffle(cards, history, random.ExpFloat64)
				if cards[0].card.Header.Content == "A" {
					firstA++
				}
			}

			frequency := float64(firstA) / trials
			if math.Abs(frequency-tt.chanceA) > 0.02 {
				t.Fatalf("first-card probability %f, want about %f",
					frequency, tt.chanceA)
			}
		})
	}
}

func TestWeightedOrderPreservesEveryOccurrence(t *testing.T) {
	var cards []practiceCard

	for i := 0; i < 30; i++ {
		card := selectionCard("Duplicate")
		card.card.Lines = []model.Line{{Original: fmt.Sprint(i)}}
		cards = append(cards, card)
	}

	random := rand.New(rand.NewPCG(15, 3))
	weightedShuffle(cards, nil, random.ExpFloat64)
	seen := map[string]bool{}

	for _, card := range cards {
		occurrence := card.card.Lines[0].Original
		if seen[occurrence] {
			t.Fatal("weighted order duplicated an occurrence")
		}

		seen[occurrence] = true
	}

	if len(seen) != 30 {
		t.Fatal("weighted order dropped occurrences with duplicate identities")
	}
}
