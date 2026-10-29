# In-memory user sessions and shuffled retries

A `Session` owns practice state independently of terminal rendering. It checks answers and builds retry rounds from missed card occurrences, while the TUI translates user input into session operations. This record describes the current implementation as of 2026-10-26.

## Starting a round

The executable requires `--deck PATH` and accepts `--limit`, defaulting to 10. File-reading errors, missing deck paths, nonpositive limits, and decks with no cards stop startup with an error.

`NewSession` copies the deck's card slice, shuffles the copy, and selects `min(limit, card count)` cards. It does not reorder the source deck. Shuffling uses `math/rand/v2`; an oversized positive limit uses all available cards.

A round tracks its selected cards, current card and answer positions, per-answer results for the current card, whether feedback is awaiting acknowledgement, and completed card occurrences that were missed. Answers are extracted from hidden line parts in line and part order.

## Answer lifecycle

1. `Submit` compares the input with the current expected answer. Both strings pass through `strings.TrimSpace`; the remaining text must match exactly. Matching is case-sensitive and preserves internal whitespace differences.
2. The result is recorded immediately. The expected answer becomes available for display, and the session waits for acknowledgement. Submission does not advance the answer position.
3. `Continue` acknowledges feedback and advances to the next answer. Previous results remain available until the card ends.
4. After the final answer is acknowledged, the session records the card as missed if any answer was incorrect, advances to the next card, and resets answer state.

`Submit` has no effect during feedback, on reading cards, or after the round is complete. `Continue` cannot skip an unanswered answer. A reading card has no hidden parts and completes with one `Continue`, without becoming missed. Operations after completion are harmless; `Card` returns an empty card.

The TUI maps Enter to submission while answering, and to `Continue` during feedback or on a reading card. It ignores ordinary typing during feedback and clears the input after continuing. The same Enter press never both submits and acknowledges an answer.

## Completion and retry rounds

A round is done when its card position reaches its card count. Its summary reports the current round's total cards and missed card occurrences, rather than a count of incorrect answers or a cumulative score.

`Repeat` starts a new round only after completion and only when there are missed cards. It uses those missed occurrences as the new card slice, shuffles them, clears the missed list, and resets card and answer state. The initial limit is not applied again. All answers are hidden again, including previously correct answers on a missed card.

A card occurrence is added once even if several answers were wrong. Identical cards appearing multiple times are tracked separately: if two occurrences were missed, both enter the retry round. No content-based deduplication occurs.

The summary initially selects `Repeat missed` when mistakes exist; the user can repeat or quit. A round without missed cards offers Enter to quit. Repeating can continue until a round has no misses or the user exits.

## Lifetime and boundaries

Session state exists only in memory. There are no accounts, saved results, scheduling intervals, or persisted learning history. Quitting discards the current state. The session contains no terminal input, colors, borders, or scrolling behavior; those belong to the [TUI](0003-ui-and-color-scheme.md).

## Implementation

- [Session and round transitions](../../internal/session/session.go)
- [Answer matching, limits, reading cards, and duplicate occurrence tests](../../internal/session/session_test.go)
- [Startup and flags](../../cmd/tui/main.go)
- [TUI input and summary handling](../../internal/tui/model.go)
