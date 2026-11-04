# Explicit deck validation

Canon D provides `canon_d --validate-deck FILE` as a CLI authoring check that
reports deck mistakes without starting the TUI or accessing practice history.
Validation is deliberately stricter than the permissive practice parser, so authors
can find mistakes without changing how existing decks behave during practice.

The design was confirmed before implementation.

## Command behavior

- Validate exactly one file. Folder validation is outside the initial scope.
- Reject combinations with `--deck`, `--decks`, an explicitly supplied `--limit`,
  or positional arguments. The default practice limit has no effect on validation.
- Do not create, read, lock, or modify practice history.
- Leave normal practice behavior unchanged, including its handling of duplicate
  identities described in [saved card statistics](0005-saved-answer-totals.md).

## Validation rules

Use the heading hierarchy, whitespace handling, fence recognition, and answer
boundaries from [deck syntax and parsing](0001-deck-syntax-and-parsing.md), with the
following stricter authoring rules.

Report these as errors:

- **Empty deck:** the file produces no usable cards.
- **Empty leaf card:** a heading has no body, or only whitespace, and no child
  headings. Bodyless section headings that group child cards are allowed.
- **Empty heading title:** the title is empty after the parser's whitespace trimming.
- **Slash in a heading:** any heading title contains `/`, including section headings
  without bodies. This keeps the full-path separator unambiguous.
- **Duplicate card path:** two actual cards have the same full hierarchical heading
  path within the file. Compare case-sensitively after the parser's title trimming,
  and report both source locations. Section headings without bodies do not
  participate in duplicate detection.
- **Unclosed fenced block:** an opening fence has no matching closing fence.
- **Empty fenced block:** a block contains no non-whitespace code lines. Check each
  block independently, even when its card also contains prose or other answers.
- **Invalid inline answer syntax:** outside fenced code, answers must use paired
  single backticks on the same line and start with a non-whitespace character.
  Reject unmatched backticks, unsupported multiple-backtick delimiters, empty or
  whitespace-only spans, and answers beginning with whitespace. Backticks within
  fenced code remain ordinary answer text.

Report nonblank content before the first recognized heading as a warning because
the parser ignores it.

Cards without hidden answers remain valid reading cards. Their presence alone is
neither an error nor a warning. An empty fenced block is still an error under its
own rule, even though the practice parser currently accepts it as reading content.

## Report and exit codes

Collect all validation findings and print a plain-text report to stdout in source
order. Findings identify the file and one-based source line numbers; duplicate
findings include both locations. Finish with counts of total cards, reading cards,
warnings, and errors. Reading cards are included in the total card count.

- Exit `0` when validation completes without errors, including warnings-only results.
- Exit `1` when validation completes with errors.
- Exit `2` for invalid command usage or file-access failures. Print operational
  failures to stderr.

The stricter rules belong to this explicit validation command. They do not make
validation a prerequisite for opening a deck in normal practice.

## Implementation

- [Parser validation](../../internal/parser/validate.go) reuses the practice parser's
  heading and fence recognition, retains source locations for diagnostics, and uses
  parsed cards for total and reading-card counts. Findings are sorted by source line.
- [The command entry point](../../cmd/tui/main.go) handles validation before opening
  history or starting the TUI, and owns report formatting and exit codes.
- [Validation tests](../../internal/parser/validate_test.go) cover authoring rules,
  hierarchy, formatting boundaries, and diagnostic locations.
  [Command tests](../../cmd/tui/main_test.go) cover output, exit codes, flag conflicts,
  and validation alongside corrupt history without modifying it or creating locks.
