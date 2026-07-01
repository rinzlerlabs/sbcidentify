.PHONY: all generate fmt vet build test check

all: check

generate:
	go generate ./boardtype/...

fmt:
	gofmt -w .

vet:
	go vet ./...

build:
	go build ./...

test:
	go test -race ./...

check: generate fmt vet build test
