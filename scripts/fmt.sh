#!/usr/bin/env bash

# This script runs various formatters for the repository.

# .proto formatter
go run github.com/bufbuild/buf/cmd/buf@v1.4.0 format -w

# go formatter
go fmt ./...