.PHONY: compile run run-tui test

compile:
	go build -o canon_d ./cmd/tui

run: compile
	./canon_d --deck ./decks/golang.md

test:
	go test ./...
	go vet ./...
