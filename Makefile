.PHONY: compile compile-cli compile-tui run run-tui test

compile: compile-cli compile-tui

compile-cli:
	go build -o canon_d ./cmd/cli

compile-tui:
	go build -o canon_d_tui ./cmd/tui

run: compile-cli
	./canon_d --deck ./decks/golang.md

run-tui: compile-tui
	./canon_d_tui --deck ./decks/golang.md

test:
	go test ./...
	go vet ./...
