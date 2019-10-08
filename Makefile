.PHONY: test ci install

all:
	go install ./cmd/mantra

test:
	go test -v ./...

ci: install all
	go test -v -race ./...

install:
	go mod download
	go mod verify
