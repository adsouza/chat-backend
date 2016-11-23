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
		log.Fatalf("Could not open connection to DB: %v.", err)
	}
	defer db.Close()
	if _, err := db.Exec(storage.PragmaCmd); err != nil {
		log.Printf("Unable to enable foreign key constraints in DB: %v.", err)
	}
	db.Exec(storage.UserTableInitCmd)
	db.Exec(storage.ConversationTableInitCmd)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Could not bind to port: %v.", err)
	}
	grpcServer := grpc.NewServer()
	ctlr := logic.NewUserController(storage.NewSQLDB(db))
	api.RegisterChatServer(grpcServer, api.NewChatServer(ctlr))
	log.Println("Chat service is now ready!")
	grpcServer.Serve(lis)
}
