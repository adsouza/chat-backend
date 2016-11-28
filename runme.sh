protoc -I ./ api/api.proto --go_out=plugins=grpc:. && \
go test storage/sqlite_test.go && \
go test logic/*_test.go && \
go run integration_demo.go && \
echo "All tests pass :-)"
go run main.go
