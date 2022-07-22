#!/bin/ash

if [[ $# -eq 1 ]]; then
        file="$1"
else
        echo "Must be called with file path as first positional argument."
        exit 1
fi

name="$(echo "$file" | awk -F . '{ print $1 }')"
ext="$(echo "$file" | awk -F . '{ print $2 }')"

fn="$name-$(date +%s).$ext"
cp "$file" "$fn"
truncate -s 0 "$file"

echo $fn

# usage:
# flock <metrics.out> ./copy-and-truncate.sh <metrics.out>