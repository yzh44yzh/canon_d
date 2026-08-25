compile:
	go fmt
	go build
	
run:
	./canon_d ./decks/golang.md
	
check:
	go test -v ./..
	go vet
