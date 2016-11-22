protoc -I api/ api/api.proto --go_out=plugins=grpc:api && \
go test storage/sqlite_test.go && \
go test logic/users_test.go && \
go run integration_test_main.go && \
echo "All tests pass :-)"
go run main.go
