BINARY := golblog

.PHONY: all build run test clean

all: build run

build:
	go build -o $(BINARY) ./cmd

run:
	go run ./cmd

test:
	go test ./...

clean:
	rm -rf out/ $(BINARY)
