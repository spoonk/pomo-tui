DEBUG_FLAGS := -gcflags="all=-N -l"

.PHONY: build run test

build:
	go build -o pomo-tui .

build-debug:
	go build $(DEBUG_FLAGS) -o pomo-tui .

run:
	make run-debug

run-debug:
	make build-debug && ./pomo-tui

clean: 
	rm -rf pomo-tui

test:
	gotestsum --format testdox ./...
