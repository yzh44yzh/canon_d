# Multi-deck learning

This document records the confirmed behavior of folder selection and multi-deck
practice.

## Input and discovery

- Add `--decks DIR` to open a deck-selection screen.
- Keep `--deck FILE` for starting practice directly with one deck.
- Reject supplying both flags together.
- List non-hidden regular `.md` files directly inside the supplied directory, without
  searching subdirectories. Sort by filename and display the `.md` suffix.

## Deck selection

- The screen title is `Select Decks`.
- Initially, no decks are selected.
- Up and Down move the cursor; Space toggles the focused usable deck.
- Enter starts learning only when at least one usable deck is selected.
- Show checkboxes, each deck's card count, and the number of selected decks.
- Empty or unreadable decks remain visible but disabled, with a reason.
- If no usable decks exist, report an explanatory error and exit.
- When the list exceeds the available height, scroll to keep the cursor visible.
  Page Up and Page Down move by a page. The title, selection count, and keyboard
  hints remain visible.
- Load deck contents when opening the selection screen. Card counts and practice
  use this same snapshot; restarting picks up subsequent file changes.

## Practice

Combine cards from the selected decks, generate a weighted random order, and apply
the existing limit once to that collection. The displayed deck name follows the
current card's source deck. Answer checking, reading cards, and shuffled retry
rounds retain their current behavior.

The [priority policy](0006-priority-weighted-card-selection.md) uses the saved
statistics snapshot for the run. The maximum old-card weight is computed only from
cards in the selected decks; unselected decks and stale history rows do not affect it.

The limit applies globally, without per-deck quotas. A selected deck can therefore
contribute no cards to the initial round. Card progress and the completion summary
describe the combined round; missed occurrences retain their source decks when
retried. The completion screen keeps its existing repeat-or-quit behavior.

## Implementation

- [Deck loading and snapshots](../../internal/decks/load.go)
- [Deck selection](../../internal/tui/picker.go)
- [Combined rounds and source-deck tracking](../../internal/session/session.go)

## Example terminal flow

For a folder containing two usable decks and one empty deck, selecting the first
two can look like this:

```text
Select Decks
Selected decks: 2
╭────────────────────────────────────────╮
│  [x] golang.md · 20 cards               │
│> [x] network.md · 15 cards              │
│  [-] empty.md · deck contains no cards  │
╰────────────────────────────────────────╯
↑/↓ move · Space toggle · Enter learn · PgUp/PgDown · Ctrl+C quit
```

With `--limit 10`, Enter starts a shuffled ten-card round. The progress header
can show `network.md · Card 1/10` followed by `golang.md · Card 2/10`, depending
on the shuffle. The card heading, answer prompts, and feedback work as before.
