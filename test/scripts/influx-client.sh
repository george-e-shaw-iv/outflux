#!/usr/bin/env bash

# Script that allows the caller to pass a flux query to the InfluxDB instance
# that runs on the client network (the one that telegraf dual-writes to by
# proxying requests).

query='from(bucket:"test") |> range(start:-5m)'
if [[ "$#" -eq 1 ]]; then
    query="$1"
fi

docker exec $(docker ps -aqf "name=test-influx-client-1") influx query "$query"