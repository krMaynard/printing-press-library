.PHONY: build test lint install clean

build:
	go build -o bin/naver-privacy-pp-cli ./cmd/naver-privacy-pp-cli

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/naver-privacy-pp-cli

clean:
	rm -rf bin/

build-mcp:
	go build -o bin/naver-privacy-pp-mcp ./cmd/naver-privacy-pp-mcp

install-mcp:
	go install ./cmd/naver-privacy-pp-mcp

build-all: build build-mcp
