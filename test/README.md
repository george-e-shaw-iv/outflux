# `outflux` Test Environment

## Pre-requisites

* [InfluxDB CLI v2](https://docs.influxdata.com/influxdb/cloud/tools/influx-cli/#install-the-influx-cli)
* [Docker](https://docs.docker.com/get-docker/)

## Running

```bash
docker compose up
```

## Interacting

The `scripts/` directory contains some handy scripts to help interact with and
visualize what is happening within the test environment:

* `dummy-write.sh`: This script takes one integer argument and writes
  that number of dummy metrics into the influx client by way of telegraf, which
  also writes it to the output file that `outflux` (client) reads from to sync
  with the server.
    * Usage: `./scripts/dummy-write.sh 10`
* `output-file.sh`: This script takes takes one optional string argument, either
  `dump` or `tail`, defaulting to `dump`. It either dumps the contents of the file
  that telegraf dual-writes to whenever it proxies a metric to the client InfluxDB
  instance, or tails the file which is a blocking operation that updates on
  changes to the file.
    * Usage: `./scripts/output-file.sh tail`
* `influx-client.sh` This script takes one optional string argument in the form of
  a [Flux](https://docs.influxdata.com/flux/v0.x/spec/) query. If not supplied,
  the query that is passed is `from(bucket:"test") |> range(start:-5m)` which will
  output the last 5 minutes of data written to the client InfluxDB instance by the
  `dummy-write.sh` script. This query is passed directly to `influx query`.
    * Usage: `./scripts/influx-client.sh 'from(bucket:"test") |> range(start:-30m)'`

## Destroying

If you want to preserve data:
```bash
docker compose down
```

If you don't want to preserve data:
```bash
docker compose down --volumes
```