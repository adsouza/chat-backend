package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/adsouza/chat-backend/api"
	"github.com/adsouza/chat-backend/logic"
	"github.com/adsouza/chat-backend/storage"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	db, err := sql.Open("sqlite3", "")
	if err != nil {
		log.Fatalf("Could not open connection to DB: %v.", err)
	}
	defer db.Close()
	if _, err := db.Exec(storage.UserTableInitCmd); err != nil {
		log.Fatalf("Unable to create new users table in test DB: %v.", err)
	}

	lis, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatalf("Could not bind to port: %v.", err)
	}
	grpcServer := grpc.NewServer()
	ctlr := logic.NewUserController(storage.NewSQLDB(db))
	api.RegisterChatServer(grpcServer, api.NewChatServer(ctlr))
	go grpcServer.Serve(lis)

	conn, err := grpc.Dial(":12345", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to server: %v.", err)
	}
	defer conn.Close()

	client := api.NewChatClient(conn)
	_, err = client.CreateUser(context.Background(), &api.CreateUserRequest{Username: "testuser1", Passphrase: "0123456789abcdef"})
	if err != nil {
		log.Fatalf("Could not create a user account: %v.", err)
	}
}
