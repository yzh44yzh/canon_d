# Validate history before practice startup

Canon D validates the entire `.canon_d.csv` before starting either single-deck
practice or the multi-deck selection TUI. Strict validation prevents corrupt or
ambiguous statistics from silently affecting practice and later saves. These
confirmed rules replace the provisional CSV validation policy in
[saved card statistics](0005-saved-answer-totals.md).

## Startup and preservation

- Validate every row, including identities for removed or renamed decks and cards.
  Historical identities need not match any current deck or card.
- A missing CSV is a valid first run with empty statistics. Create history only
  when practice results are saved. An existing empty or blank-only file is invalid;
  a header-only file is valid.
- Stop at the first error, before showing any TUI. Print the file path, physical
  source line, field name where applicable, and reason to stderr; exit `1`.
- Leave invalid CSV bytes untouched. Do not discard rows, reset totals, merge
  duplicates, or repair the file automatically. Release the history lock on failure;
  the lock file may remain, as during normal operation.
- `--validate-deck` remains independent of history and its lock.

## CSV structure and identities

Require this exact decoded header, in this order, and five fields per data record:

```csv
deck,card_header,shows,correct_answers,incorrect_answers
```

Use standard CSV quoting and escaping. Accept blank lines, LF or CRLF, and an
optional final newline. Header spaces, a byte-order mark, and comments are not
special syntax and do not bypass the exact header or record requirements.

Both identity fields must contain non-whitespace text, valid UTF-8, and no control
characters. The `deck` field must be a filename: reject `/`, `\`, `.` and `..`;
no `.md` extension is required. The `card_header` field stores the full hierarchy,
so separators such as `Go / Basics` are valid. Quoted commas and quotes remain
valid, but embedded newlines and other control characters are invalid identities.

Compare `(deck, card_header)` pairs exactly, preserving case and whitespace.
Reject repeated pairs even when their counters match. Report the duplicate's
location and the first occurrence's location.

## Counters

Require decimal unsigned 64-bit integers in the inclusive range
`0` through `18446744073709551615`. Leading zeros are allowed. Reject signs,
surrounding whitespace, empty values, fractions, and overflow.

Positive correct or incorrect answer totals require `shows > 0`. All-zero rows
and shows without answers are valid. Do not impose an answers-per-show limit:
one card can contain multiple answers. Accept maximum counter values at startup;
retain overflow protection when saving additional results.

## Implementation

[CSV loading](../../internal/history/csv.go) validates decoded identities, counters,
and duplicate keys, using physical CSV field positions for diagnostics.
[The history store](../../internal/history/store.go) acquires the folder lock before
loading and releases it on failure. [Practice startup](../../cmd/tui/main.go) opens
the history before loading decks or starting the TUI.
