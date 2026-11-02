package session

import (
	"cmp"
	"slices"

	"com.github.yzh44yzh/canon_d/internal/model"
)

// initialWeights uses only saved totals for occurrences in the eligible pool.
func initialWeights(
	cards []practiceCard,
	history map[model.CardIdentity]model.CardStatistics,
) []float64 {
	weights := make([]float64, len(cards))
	maxOld := 0.0

	for i, card := range cards {
		id := model.CardIdentity{Deck: card.deck, Header: card.card.Header.Show()}
		total := history[id]
		if total.CorrectAnswers == 0 && total.IncorrectAnswers == 0 {
			continue
		}

		// Convert before adding so maximum uint64 counters cannot overflow.
		errors := float64(total.IncorrectAnswers) + 1
		correct := float64(total.CorrectAnswers) + 1
		weight := errors / correct
		weights[i] = weight
		maxOld = max(maxOld, weight)
	}

	newWeight := 1.0
	if maxOld > 0 {
		newWeight = 2 * maxOld
	}

	for i, weight := range weights {
		if weight == 0 {
			weights[i] = newWeight
		}
	}

	return weights
}

// weightedShuffle assigns independent exponential arrival times. A card's rate is
// its weight, making its chance of being next proportional to its remaining weight.
// Each occurrence gets one arrival time and therefore appears exactly once.
func weightedShuffle(
	cards []practiceCard,
	history map[model.CardIdentity]model.CardStatistics,
	exponential func() float64,
) {
	type arrival struct {
		card practiceCard
		time float64
	}

	weights := initialWeights(cards, history)
	arrivals := make([]arrival, len(cards))

	for i, card := range cards {
		arrivals[i] = arrival{card: card, time: exponential() / weights[i]}
	}

	slices.SortFunc(arrivals, func(a, b arrival) int { return cmp.Compare(a.time, b.time) })

	for i, entry := range arrivals {
		cards[i] = entry.card
	}
}
