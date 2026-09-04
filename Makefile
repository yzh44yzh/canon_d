compile:
	go build -o canon_d ./cli/main.go
	
run:
	./canon_d ./decks/golang.md
	
test:
	go test ./...
	go vet ./...
