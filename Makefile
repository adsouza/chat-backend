.PHONY: compile
GOPATH?=$(HOME)/go
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC := $(shell which protoc)
# If protoc isn't on the path, set it to a target that's never up to date, so
# the install command always runs.
ifeq ($(PROTOC),)
    PROTOC = must-rebuild
endif

# Figure out which machine we're running on.
UNAME := $(shell uname)

all: test

$(PROTOC):
# Run the right installation command for the operating system.
ifeq ($(UNAME), Darwin)
	brew install protobuf
endif
ifeq ($(UNAME), Linux)
	sudo apt-get install protobuf-compiler
endif
# You can add instructions for other operating systems here, or use different
# branching logic as appropriate.

# If $GOPATH/bin/protoc-gen-go does not exist, we'll run this command to install
# it.
$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

api/api.pb.go: api/api.proto | $(PROTOC_GEN_GO) $(PROTOC)
	protoc -I ./ --go_out=plugins=grpc:. api/api.proto

# This is a "phony" target - an alias for the above command, so "make compile"
# still works.
compile: api/api.pb.go

clean:
	rm api/api.pb.go
	go clean
	rm chat.db

test: compile
	go test storage/sqlite_test.go
	go test logic/*_test.go
	go run integration_demo.go
