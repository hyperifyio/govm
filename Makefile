.PHONY: build run clean tidy certs

GOVM_SOURCES := $(shell find ./*.go ./cmd ./internal -type f -iname '*.go' ! -iname '*_test.go')

all: build certs

tidy:
	go mod tidy

build: govm

govm: $(GOVM_SOURCES) Makefile
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o govm ./cmd/govm

test: Makefile
	go test -v ./...

certs:
	make -C certs

clean:
	rm -f ./govm
