BINARY := golblog

.PHONY: all build run test clean

all: build run

build:
	go build -o $(BINARY) .

run:
	go run main.go

test:
	go test ./...

clean:
	rm -rf out/ $(BINARY)
