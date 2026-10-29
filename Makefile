.PHONY: compile run check test

compile:
	go build -o canon_d ./cmd/tui

run: compile
	./canon_d --deck ./decks/golang.md

check:
	go run ./scripts/check

test: check
	go test ./...
	go vet ./...
