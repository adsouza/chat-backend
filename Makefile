.PHONY: compile
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go

# If $GOPATH/bin/protoc-gen-go does not exist, we'll run this command to install
# it.
$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

api/api.pb.go: api/api.proto | $(PROTOC_GEN_GO)
	protoc -I ./ --go_out=plugins=grpc:. api/api.proto

# This is a "phony" target - an alias for the above command, so "make compile"
# still works.
compile: api/api.pb.go

clean:
	rm api/api.pb.go
