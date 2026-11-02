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
- For now, invalid existing CSV data stops startup. A failed save stops practice
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
- For now, loading requires the exact CSV header, five fields per row, unsigned
  integer counters, and unique card identities. Malformed data stops startup without
  rewriting the file. These checks are provisional pending the deferred interview.
- [The TUI](../../internal/tui/model.go) records visible presentations and requests
  saves at round completion and on Ctrl+C. The CLI owns the history store and reports
  failures after the terminal UI exits.

For example, after a card with two answers is shown and receives one correct and
one incorrect answer, its row can be:

```csv
deck,card_header,shows,correct_answers,incorrect_answers
database.md,SQL / Isolation,1,1,1
```

A retry with both answers correct changes those totals to `2,3,1`.

## Deferred interviews

The user explicitly requested later design interviews about deck validation rules,
invalid-deck behavior, and invalid-CSV behavior. The provisional decisions above
apply until those interviews revise them; those topics are not part of the current
interview's remaining scope.
