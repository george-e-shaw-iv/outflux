# Contributing to `outflux`

## Pre-requisites

* [`asdf`](https://asdf-vm.com/guide/getting-started.html)
    * Once you have this installed run `./scripts/asdf.sh` which will add all of the plugins
      you need and then run `asdf install` (this script needs to be invoked in the root of this
      repository).
    * If you do not want to install `asdf` for whatever reason, you can manually install the
      versions (at least match the major version on each) of each tool in `.tool-versions`.
* If you're going to be running the test suite, you'll need all of the pre-requisites it requires
  which can be found [here](test/README.md#pre-requisites).

## Building

```bash
make build
```

## Generating Protobuf Files

```bash
make generate
```