.PHONY: build run clean tidy certs config

GOVM_SOURCES := $(shell find ./*.go ./cmd ./internal -type f -iname '*.go' ! -iname '*_test.go')

all: build config

tidy:
	go mod tidy

build: govm

config: config.yml certs

config.yml:
	@echo "servers: []" > $@

govm: $(GOVM_SOURCES) Makefile config
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o govm ./cmd/govm

test: Makefile
	go test -v ./...

certs:
	make -C certs

clean:
	rm -f ./govm
