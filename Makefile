.PHONY: all test ci install

install:
	go mod download
	go mod verify

all:
	go install ./cmd/mantra

test:
	go test -v ./...

ci:
	go mod download
	go test -v -race
