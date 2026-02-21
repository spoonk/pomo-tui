.PHONY: build run test

build:
	go build -o pomo-tui .

run:
	go run .

test:
	gotestsum --format testdox ./...
