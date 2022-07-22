#!/usr/bin/env bash

if [[ "$(go version | grep 1.18)" ]]; then
    echo "grpcui does not work for go1.18 right now"
    exit 0
fi

GRPCUI="$(go env | grep GOPATH | awk -F \" '{print $2}')/bin/grpcui"

if [[ ! -f "$grpcui" ]]; then
    go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
fi

"$GRPCUI" localhost:8000