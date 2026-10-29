# Deck syntax and formatting-preserving parsing

Canon D stores decks as Markdown files and parses them into cards with visible and hidden text parts. The parser preserves source formatting so the TUI can hide answers without losing the surrounding context. This record describes the current implementation as of 2026-10-26.

## Deck and card structure

A `Deck` has a header and an ordered slice of cards. The entry point uses the deck file's base name as the deck header. Each `Card` has a hierarchical `CardHeader` and a slice of `Line` values. A line stores its original text and ordered `LinePart` values; each part contains text and a visibility flag.

Outside fenced code, a line starting with `#` in the first column starts a card and ends the previous one. The number of leading `#` characters determines the heading level. The remaining heading text is trimmed. A space after `#` is optional, and there is no six-level restriction. Indented headings and Setext headings are not recognized.

A heading's parent is the nearest preceding heading with a smaller level. `CardHeader.Show` joins ancestor titles with ` / `. A heading with no body or only whitespace creates no card, but can still be a parent. Body text before the first recognized heading creates no card. A nonempty body without hidden answers creates a reading card.

## Answer syntax

- Outside fences, single-backtick spans mark hidden answers. The match must start with a non-whitespace character and end with a backtick on the same line. Multiple spans produce separate answers in source order. This is a small regular-expression parser, without Markdown escape or variable-length inline-delimiter handling.
- A fence opens when a trimmed line begins with at least three consecutive backticks or tildes. The opening run determines the delimiter character and minimum closing length; trailing opening text, such as a language tag, is retained.
- A closing fence consists only of the opening delimiter character, repeated at least as many times, after trimming surrounding whitespace. Shorter markers remain code content.
- Each non-whitespace code line becomes one hidden answer. Leading spaces and tabs remain visible indentation; the rest is the answer, including any trailing whitespace. Inline spans inside code are not parsed separately.
- Fence markers and whitespace-only lines stay visible and require no answer. Headings inside a fence are code, not card boundaries. An unclosed fence extends to the end of the input; the parser returns no syntax error.

## Formatting and display boundary

`MakeDeck` normalizes CRLF to LF and removes one final newline before splitting the input. Within card bodies it retains indentation, blank lines, fence markers, and inline backticks. Concatenating a line's parts reconstructs its `Original` value. A body containing fence markers is nonempty even if the fenced block has no answers.

The TUI separately removes backticks from display, expands tabs, filters control characters, wraps text, and inserts numbered placeholders. It does not change the stored answer text. Answer checking belongs to the [session](0002-user-session.md).

## Example

````markdown
# Networking
## DNS
Resolve names with `DNS`.

```sh
  dig example.org
```
## Reading
Review the resolver configuration.
````

This produces two cards, headed `Networking / DNS` and `Networking / Reading`. The first has answers `DNS` and `dig example.org`; the second is a reading card. The parent heading has no body and produces no card.

## Implementation

- [Deck, card, header, and line types](../../internal/model/model.go)
- [Parser](../../internal/parser/parser.go)
- [Parser tests](../../internal/parser/parser_test.go)
- [File loading and deck naming](../../cmd/tui/main.go)
