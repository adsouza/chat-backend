# Chat Back-end

This is a Go implementation for a very simple chat backend. It uses SQLite3 for storage and gRPC for the API.

## Installation

1. Download the OS-specific Protocol Buffers compiler from https://github.com/google/protobuf/releases and unpack it into /usr/local
1. `go get -u github.com/golang/protobuf/{proto,protoc-gen-go}`
1. `export PATH=$PATH:$GOPATH/bin`
1. `go get github.com/adsouza/chat-backend`

## Usage

`./runme.sh` will generate gRPC bindings, run all tests and start the server.

`./clean.sh` will delete the generated gRPC bindings and the SQLite3 DB file.

## Storage

The storage module has a SQL implementation that has been tested with SQLite3.

## Logic

The logic module is pure Go and relies upon a storage interface for which an implementation is available in the corresponding module.

## API

The API is built using gRPC and relies upon a controller interface that is implemented by the logic module.

