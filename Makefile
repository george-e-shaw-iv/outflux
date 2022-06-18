PHONY: build
build: fmt generate
	go build -o bin/ ./cmd/...

PHONY: fmt
fmt:
	./scripts/fmt.sh

PHONY: generate
generate:
	go generate ./...