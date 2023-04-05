SHELL=/bin/bash

.PHONY: build
build:
	CGO_ENABLED=0 go build -a -o ./build/server cmd/main.go
	CGO_ENABLED=0 go build -a -o ./build/tracker tracker/main.go

.PHONY: protogen
protogen:
	bash protocgen.sh

.PHONY: test
test:
	go test -v -count=1 ./...

.PHONY: run-server
run-server:
	go run cmd/main.go

.PHONY: run-tracker
run-tracker:
	go run tracker/main.go

.PHONY: lint
lint:
	golangci-lint run ./... --fix