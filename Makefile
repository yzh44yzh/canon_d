compile:
	go build -o canon_d ./cli/main.go
	
run:
	./canon_d --deck ./decks/golang.md
	
test:
	go test ./...
	go vet ./...
