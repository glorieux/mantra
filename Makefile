.PHONY: test ci install

all:
	go install ./cmd/mantra

test:
	go test -v ./...

test-race:
	go test -v -race ./...

ci: install all test-race

pre-commit: test-race
	go mod tidy
	go mod verify
	go vet ./...

install:
	go mod download
	go mod verify
