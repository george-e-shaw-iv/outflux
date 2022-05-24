#!/usr/bin/env bash

# This script generates protobuf files.

if [[ $(pwd) != */internal/grpc ]]; then
    echo "This script should only be ran via make generate (pwd:$(pwd))."
    exit 1
fi

# Install the plugins we need for protoc.
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# This is where the plugins will be installed, the path could definitely change
# since we're using asdf to manage go versions.
gobin="$(go env GOPATH)/bin"

# Print versions
"$gobin/protoc-gen-go" --version
"$gobin/protoc-gen-go-grpc" -version

# Run protobuf codegen.
protoc --plugin=protoc-gen-go="$gobin/protoc-gen-go" \
    --plugin=protoc-gen-go-grpc="$gobin/protoc-gen-go-grpc" \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api.proto