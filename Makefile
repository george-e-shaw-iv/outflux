PHONY: build
build: fmt generate
	go build -o bin/ ./cmd/outflux

PHONY: fmt
fmt:
	./scripts/fmt.sh

PHONY: generate
generate:
	go generate ./...