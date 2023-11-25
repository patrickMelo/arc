# ARC

This is a virtual machine based simple Redis-like in memory database created for a challenge.

## How It Works

* An HTTP server is spawned to listen for connections and commands.
* Commands are parsed and executed in a virtual machine environment.
* Available commands are defined in runtime libraries, that can be easily expanded or exchanged.
* For each command a result set (or error message) is returned and then sent back to the requesting client.

## Compiling

Just run `go build` in the `source/` folder (optionally add `-o ../binaries/arc` to output the program on the binaries folder).

A ready to run linux-x64 binary is available at `binaries/arc.x64`.

## Running

ARC can be run in three modes: client, server and standalone; just run `arc [mode]`.

* `client`: runs an interactive shell client where user can issue database commands (currently it connects only to `localhost:8080`).
* `server`: runs a HTTP server that accepts command via the `cmd` query parameter or a `REST` request (currently it only runs on `localhost:8080`).
* `standalone`: runs an interactive shell that executs commands in memory, no server or client is spawned.

## Unit Tests

Unit tests are available for the `database` package. Just run `go test` in the `source/database/` folder.

## HTTP Tests

Test files for `cmd` and `REST` requests are available in the `tests/` folder. You can use the VSCode extension [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) for testing.

