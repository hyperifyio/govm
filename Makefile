.PHONY: build run clean tidy certs config update update-frontend

GOVM_SOURCES := $(shell find ./*.go ./cmd ./internal -type f -iname '*.go' ! -iname '*_test.go')

all: build config

tidy:
	go mod tidy

build: govm

docs: api.html api.md

openapi.json: $(GOVM_SOURCES) Makefile
	mkdir -p ./tmp
	curl --insecure https://localhost:3001/documentation/json -o ./tmp/openapi.json
	mv -f ./tmp/openapi.json openapi.json

# API docs as HTML from OpenAPI specification
api.html: openapi.json Makefile
	mkdir -p ./tmp
	swagger-codegen generate -i openapi.json -l html -o ./tmp
	mv -f ./tmp/index.html api.html

# Markdown version for the API docs
api.md: api.html Makefile
	mkdir -p ./tmp
	pandoc ./api.html -f html -t markdown -o ./tmp/api.md
	mv -f ./tmp/api.md api.md

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

clean-docs:
	rm -f ./api.html ./api.md ./openapi.html ./openapi.json ./tmp/.swagger-codegen-ignore ./tmp/.swagger-codegen/VERSION
	test -e ./tmp/.swagger-codegen && rmdir ./tmp/.swagger-codegen || true
	test -e ./tmp && rmdir ./tmp || true
