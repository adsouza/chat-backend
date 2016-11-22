# Chat Back-end

This is a Go implementation for a very simple chat backend. It uses SQLite3 for storage and gRPC for the API.

## Usage

`./runme.sh` will generate gRPC bindings, run all tests and start the server.

`./clean.sh` will delete the generated gRPC bindings and the SQLite3 DB file.

## Storage

The storage module has a SQL implementation that has been tested with SQLite3.

## Logic

The logic module is pure Go and defines a storage interface for which an implementation is available in the corresponding module.

## API

The API is built using gRPC and defines a controller interface that is implemented by the logic module.

