.PHONY: all test build run vet critic lint

all: test build

test:
	go test ./...

build:
	go build -o wtester .

run:
	go run .

vet:
	go vet ./...

critic:
	gocritic check ./...

sec:
	gosec ./...

lint: test vet critic sec