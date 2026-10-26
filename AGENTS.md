# Repository Guidelines

## Project Structure & Module Organization

Canon D is a Go terminal flashcard trainer. The entry point is `cmd/tui/main.go`. `internal/model.go` defines decks and cards. `internal/parser.go` parses Markdown while preserving formatting. `internal/session.go` owns TUI answer checking and shuffled retries. Bubble Tea rendering and input live in `internal/tui/`. Keep domain behavior separate from rendering.

`decks/` contains Markdown learning material. `README.md` and `README.ru.md` provide English and Russian documentation; update both when changing documented behavior.

## Build, Test, and Development Commands

Use Go 1.25.0 or later. Run commands from the repository root:

- `make compile`: build `canon_d`.
- `make run`: build and run the TUI with `decks/golang.md`.
- `make test`: run all Go tests, then `go vet ./...`.
- `go run ./cmd/tui --deck ./decks/network.md --limit 5`: run the TUI directly.
- `go test ./internal -run TestPreserveFormatting -v`: run a focused test.
- `gofmt -w cmd/tui/*.go internal/*.go internal/tui/*.go`: format Go source.

The TUI validates limits and empty decks.

## Coding Style & Naming Conventions

Use `gofmt` and tabs. Use MixedCaps for exported identifiers, mixedCaps for unexported identifiers, and lowercase package names. Follow existing error-return patterns. The static check is `go vet`; there is no separate linter configuration.

## Testing Guidelines

Use Go's `testing` package and `github.com/google/go-cmp/cmp` for structured comparisons. Place tests in adjacent `*_test.go` files, with `TestXxx` names and table-driven subtests where useful. Cover formatting preservation, answer matching, retry rounds, and TUI transitions when changing those behaviors. No coverage threshold is configured. Run `make test` before submitting code changes.

## Deck Format

Markdown headings delimit cards and establish their hierarchy. Inline backticks mark hidden answers; nonempty fenced-code lines become individual answers. Empty cards are skipped. The parser retains whitespace and fences and treats headings inside fences as code.

## Commit & Pull Request Guidelines

History uses short, imperative subjects such as `shuffle Deck and limit cards` and `move modules to internal package`, without a mandatory prefix. Keep commits focused. Pull requests should explain the behavior change, link a relevant issue when available, and list validation commands and results. Include a sample terminal session for changes to prompts or answer feedback.
