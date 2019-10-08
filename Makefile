.PHONY: test ci install

all:
	go install ./cmd/mantra

test:
	go test -v ./...

ci:
	go mod download
	go test -v -race ./...

install:
	go mod download
	go mod verify
