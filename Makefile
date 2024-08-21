.PHONY: build run clean tidy certs config update update-frontend

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

update: update-frontend

update-frontend:
	cd internal/frontend-govm/ && git pull

certs:
	make -C certs

clean:
	rm -f ./govm
