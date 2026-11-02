# Priority-weighted card selection

Initial practice rounds use saved card
statistics to favor new cards and cards with more incorrect answers relative to
correct answers, while retaining random selection.

## Confirmed decisions

- A card is new when it has no history row, or when both its correct and incorrect
  answer totals are zero. A recorded show alone does not make a card old.
- Reading cards retain new-card priority indefinitely because they never accumulate
  correct or incorrect answer totals. Shows do not affect priority.
- For old cards, use weight `(incorrect_answers + 1) / (correct_answers + 1)`.
  This keeps all weights positive and handles zero correct answers. The ratio is
  based on cumulative saved totals, including earlier retry rounds.
- New cards have the highest individual weight. Old cards remain eligible even
  when there are enough new cards to fill the limit; selection is probabilistic.
- Each new card receives twice the highest old-card weight in the eligible pool
  from the selected decks. If all eligible cards are new, give them equal weight.
  Compute weights once before generating the initial order; do not recalculate
  them as cards are drawn or answered.
- Generate a weighted random order of all eligible card occurrences, then apply
  the existing global limit. Weight therefore affects both selection and order.
  Each occurrence can appear at most once in the initial round.
- Retry rounds remain unchanged: include every missed occurrence and shuffle
  uniformly, without reapplying the initial limit.
- Apply this policy to both `--deck` and `--decks`. Use the saved statistics loaded
  for the run. Duplicate identities inherit the same weight while remaining
  separate selectable occurrences, as in the existing practice model.

## Implementation

The design was confirmed before implementation.

- [The history store](../../internal/history/store.go) exposes a copy of the saved
  statistics loaded at startup. This snapshot excludes current-run results, even
  after saves. It is passed to both the direct session constructor and the deck
  picker, which creates the session after selection.
- [Initial selection](../../internal/session/selection.go) computes positive weights
  once for the eligible pool. Each occurrence receives an independent exponential
  arrival time divided by its weight; sorting arrival times produces weighted
  sampling without replacement. The session then applies its global limit.
- Convert saved counters to floating point before adding one, so maximum integer
  counters cannot wrap. The entire allowed counter range gives finite positive
  weights.
- Saved statistics are selection inputs, separate from the current-run totals
  accumulated for persistence. Selection never mutates decks or saved history.
- Retry ordering continues to use the existing uniform shuffle.
