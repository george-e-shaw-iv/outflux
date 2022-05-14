# outflux

[![Go Reference](https://pkg.go.dev/badge/github.com/george-e-shaw-iv/outflux.svg)](https://pkg.go.dev/github.com/george-e-shaw-iv/outflux)

`outflux` is a smart remote [InfluxDB](https://www.influxdata.com/) syncer that uses
client/server architecture written in Go.

## How it Works

The client connects to the server using two pieces of information:

1. `outflux` server address.
2. The credentials to connect to the same backing InfluxDB instance that the server is
   using.

Once connected to the server, the way that the server determines whether or not this
client is valid is by attempting to use the InfluxDB credentials sent to it to connect
to the database that the server is using. This prevents the need for a second
authentication mechanism past what the database already provides.

The client sets the syncing rules, which can be any of the following:

* Streaming
  * This uses an HTTP connection that is upgraded to a websocket to do real-time syncing
    whenever the client writes data to its own InfluxDB instance.
* Time-based Interval
  * On an interval specified by the client the client will aggregate data written to its
  own InfluxDB instance over a time based interval and send it at the end of the interval
  to the `outflux` server.
* Size-based interval
  * The client will aggregate data written to its own InfluxDB instance until the amount
  of new data queued to be sent to the server has been hit. Once the size has been hit,
  it sends all of the data to the `outflux` server.
* On-demand
  * This can be used in conjunction with any of the other syncing rules. It essentially
  acts as a "flush" of any data that hasn't already been written to the `outflux` server
  from the client.

The manner in which the syncing is actually accomplished is by clever use of dual-writing
the data that the client is storing in its own InfluxDB instance to a file using
[telegraf](https://www.influxdata.com/time-series-platform/telegraf/). This data is
written in [InfluxDB Line Protocol](https://docs.influxdata.com/influxdb/v2.2/reference/syntax/line-protocol/)
format and the client is configured to pull data it needs to sync to the server from this
file. When the client reads from and successfully syncs using the file, it will adjust it
accordingly by removing those now-synced lines.
