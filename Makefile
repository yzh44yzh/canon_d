compile:
	go build -o canon_d ./cli/main.go
	
run:
	./canon_d ./decks/golang.md
	
check:
	go test -v ./...
	go vet
