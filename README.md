# Canon D

> The goal of learning is not to acquire or accumulate knowledge, but to develop the ability to act with understanding.
> Pyotr Yakovlevich Galperin.

> Conceptual knowledge is not enough. You need conviction that comes from your own experience.
> Maxim Dorofeev.

A terminal tool for active [spaced repetition](https://en.wikipedia.org/wiki/Spaced_repetition).

Most spaced repetition systems use "passive" flashcards. This means you simply read a card and then decide how well you remember the information on it. In fact, you need to make an effort to recall it, and that effort is an essential part of spaced repetition.

However, the brain tries to avoid this effort, and "passive" flashcards can easily become a formality with no real benefit. An "active" flashcard requires an action—in this case, typing the missing text. This forces the brain to make an effort and recall the information.


## Meaning of the Name

> Canon (art): a set of norms and rules in art.

> Canon (science): a set of fundamental laws, norms, and methods in scientific knowledge.

Canon D -- Canon for Development -- stands for a programmer's fundamental skills. (_Perhaps this name is too ambitious for the current version. But I plan to develop the application into something more_).

Canon D (Canon in D major) is also a well-known piece by the German composer Johann Pachelbel. Thousands of different covers have been made of it. For example, [this rock version](https://www.youtube.com/watch?v=XZgiNnGB8m4).


## Why Use This When AI Knows Everything, Remembers Everything, and Can Give You All the Answers?

We train our brains much like athletes train their muscles -- through exercises that create an artificial workload. Athletes run and lift weights. Knowledge workers memorize facts and solve problems.

Having a car does not replace going for a run. Having movers does not replace going to the gym. AI does not remove the need to train your natural intelligence.


## Usage

Build with Go 1.25 or later:

```sh
make compile
```

Run with a deck:

```sh
./canon_d --deck ./decks/golang.md --limit 10
```

Or select several decks from a folder:

```sh
./canon_d --decks ./decks --limit 10
```

Use either `--deck` or `--decks`. The folder view lists non-hidden regular `.md` files
in filename order, without searching subfolders. Use Up/Down to move, Space to select
or deselect decks, and Enter to learn. No decks are selected initially. Empty or
unreadable decks are disabled with a reason; a folder without usable decks reports
an error. Card counts and practice use the contents loaded when the screen opens.
Restart to pick up file changes.

The selected cards are combined and shuffled before applying `--limit` (default 10).
The limit must be positive; an oversized limit uses all available cards. There is no
per-deck quota, so a selected deck may contribute no cards. The name above each card
identifies its source deck, including during retries. Completion and retry controls
are unchanged.

Use a terminal of at least 20 columns and 11 rows. A horizontal rule separates the fixed card header
from its scrollable content. Use PgUp/PgDown to scroll.
In deck selection, PgUp/PgDown move by a page; the title, selected-deck count, and
keyboard hints remain visible. Ctrl+C quits from any screen.

Practice statistics are saved automatically to `.canon_d.csv` in the deck folder
(the parent folder with `--deck`). Each row contains:

```csv
deck,card_header,shows,correct_answers,incorrect_answers
```

Totals accumulate across runs and retries. A show counts once when a card becomes
visible; redrawing, scrolling, and resizing do not count again. Each submitted
hidden answer adds to the correct or incorrect total. Reading cards add shows only.
Rows are identified by deck filename and full card header; duplicate headers in a
deck share totals, and renaming a deck or header starts a new history entry.

Statistics are saved after each round and on Ctrl+C, including submitted answers
on unfinished cards. A forced termination can lose changes since the last save.
Invalid CSV data prevents startup; a failed save exits with an error and preserves
the previously saved file. One run at a time may use a deck folder's history.
The hidden `.canon_d.csv.lock` file remains beside the CSV; its presence does not
mean a run is active. Both files are ignored by Git in this repository.

You can create your own flashcards; see [decks/network.md](./decks/network.md) for an example.


## Documentation

- [Deck and card syntax, parsing](./docs/adr/0001-deck-syntax-and-parsing.md)
- [User sessions, answer checking, and retries](./docs/adr/0002-user-session.md)
- [UI layout, controls, and color scheme](./docs/adr/0003-ui-and-color-scheme.md)
- [Multi-deck learning](./docs/adr/0004-multi-deck-learning.md)
- [Saved card statistics](./docs/adr/0005-saved-answer-totals.md)
- [Domain glossary](./GLOSSARY.md)


## Future Plans

Planned improvements include spaced repetition scheduling, persistent learning progress, and reports.

There may also be a multi-user web version, a mobile version, and an Obsidian plugin.
