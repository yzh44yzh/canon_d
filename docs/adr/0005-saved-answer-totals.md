# Saved card statistics

Canon D keeps cumulative card statistics
per card so practice results survive between runs. Results belong to the deck
folder, allowing single-deck and multi-deck practice to share history.

## Confirmed decisions

- Store results in `.canon_d.csv` in the folder supplied to `--decks`, or the
  parent folder of the file supplied to `--deck`.
- Create the file if absent; otherwise update its existing data.
- Keep one row per card identity: source deck filename plus full hierarchical
  card header. Renaming either starts a new history entry.
- For now, duplicate identities share one statistics row while remaining separate
  occurrences during practice.
- Store the number of shows in addition to correct and incorrect answer totals.
  Create a row when the card is first shown, including reading cards. Preserve
  existing rows for cards that have been removed or renamed.
- Count a show once when a card first becomes visible during practice. A retry
  presentation counts again. Redrawing, scrolling, resizing, revealing answers,
  and acknowledging feedback do not add shows. An undersized terminal delays the
  show until the card becomes visible. Quitting immediately after presentation
  retains that show at the next save.
- Accumulate correct and incorrect submitted answers across runs and retry
  rounds. A card with two correct answers and one incorrect answer adds two to
  its correct total and one to its incorrect total. Acknowledging feedback
  does not count as an answer.
- Save after each completed round, including retries, and on Ctrl+C. Early exit
  includes submitted answers on an unfinished card. Forced termination can lose
  changes since the last save.
- Invalid existing CSV data stops startup. A failed save stops practice
  with a clear error and preserves the previously saved file. Successful saves
  replace the file atomically.
- Refuse a second concurrent run using the same history file. Runs using different
  deck folders can proceed independently.

## CSV format

Use one header row and these five columns:

```csv
deck,card_header,shows,correct_answers,incorrect_answers
```

The deck is the filename including its extension. The card header is the full
hierarchical header used for card identity. Counters are cumulative nonnegative
integers. CSV quoting preserves commas, quotes, and other text in names and headers.

## Implementation

The design was confirmed before implementation.

- [Session statistics](../../internal/session/statistics.go) track shows and answer
  totals without accessing files.
- [History storage](../../internal/history/store.go) holds an exclusive lock on
  `.canon_d.csv.lock` for the duration of the run. The lock file remains in place;
  closing its handle releases the operating-system lock.
- [CSV storage](../../internal/history/csv.go) reads existing totals and writes rows
  sorted by deck and header. Each save combines original totals with cumulative
  results from the current run, so repeated checkpoints do not double-count.
  Unchanged results do not rewrite the file. A run with no shows creates no CSV.
- A temporary file in the same folder is written, flushed, synced, and closed before
  replacing the CSV. Failed saves preserve the previously saved CSV and clean up
  the temporary file. Counter overflow is an error.
- Before creating the temporary file, validate every merged row with the same
  identity and answer/show consistency rules used when loading. Reject the entire
  save on the first invalid row in sorted deck/header order. Preserve identities
  exactly, all previous file bytes, and the startup snapshot. A successful save
  must remain readable by the history store. The
  [history validation record](0008-history-csv-validation.md) defines save diagnostics
  and the failure policy for permissively parsed identities that cannot be stored.
- Loading requires the exact CSV header, five fields per row, unsigned integer
  counters, and unique card identities. Malformed data stops startup without
  rewriting the file. [History CSV validation](0008-history-csv-validation.md)
  defines the confirmed identity rules, counter consistency, and diagnostics.
- [The terminal adapter](../../internal/tui/model.go) reports display eligibility
  and translates keys into semantic practice actions. A private
  [practice coordinator](../../internal/tui/practice.go) sequences session transitions,
  records eligible presentations, and invokes the injected save function at round
  completion and on Ctrl+C. Checkpoints are synchronous: a completed-round action
  returns only after saving, before the summary is shown. The coordinator retains
  the first save failure and blocks subsequent actions, shows, and saves.
- The CLI still owns history opening, validation, and lock release, and reports
  failures after the terminal UI exits. Deck discovery and selection retain their
  existing owners. The coordinator starts with an existing session and performs no
  file access itself.
- [Coordinator tests](../../internal/tui/practice_test.go) exercise lifecycle decisions
  without terminal messages. [TUI statistics tests](../../internal/tui/statistics_test.go)
  retain visibility and key-mapping checks plus a real-history partial-card save.

For example, after a card with two answers is shown and receives one correct and
one incorrect answer, its row can be:

```csv
deck,card_header,shows,correct_answers,incorrect_answers
database.md,SQL / Isolation,1,1,1
```

A retry with both answers correct changes those totals to `2,3,1`.

## Deferred interviews

Explicit validation is recorded in [deck validation](0007-deck-validation.md), and
startup history checks are recorded in
[history CSV validation](0008-history-csv-validation.md). Normal practice continues
to use permissive deck parsing; changes to invalid-deck behavior during practice
remain outside these decisions.
