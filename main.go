package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/adsouza/chat-backend/api"
	"github.com/adsouza/chat-backend/logic"
	"github.com/adsouza/chat-backend/storage"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

func main() {
	dsn := flag.String("dsn", "chat.db", "Data Source Name to use for storage layer.")
	port := flag.Uint("port", 12345, "Port number on which to listen for incoming connections.")
	flag.Parse()

	db, err := sql.Open("sqlite3", *dsn)
	if err != nil {
		log.Fatalf("Could not open connection to DB: %v", err)
	}
	defer db.Close()
	// Blindly try to create the users table because if it already exists then this will safely fail.
	db.Exec(storage.UserTableInitCmd)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Could not bind to port: %v", err)
	}
	grpcServer := grpc.NewServer()
	ctlr := logic.NewUserController(storage.NewSQLDB(db))
	api.RegisterChatServer(grpcServer, api.NewChatServer(ctlr))
	log.Println("Chat service is now ready!")
	grpcServer.Serve(lis)
}
