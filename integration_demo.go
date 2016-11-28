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
	if _, err := db.Exec(storage.PragmaCmd); err != nil {
		log.Printf("Unable to enable foreign key constraints in DB: %v.", err)
	}
	if _, err := db.Exec(storage.UserTableInitCmd); err != nil {
		log.Fatalf("Unable to create new users table in test DB: %v.", err)
	}
	if _, err := db.Exec(storage.ConversationTableInitCmd); err != nil {
		log.Fatalf("Unable to create new conversations table in test DB: %v.", err)
	}

	lis, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatalf("Could not bind to port: %v.", err)
	}
	grpcServer := grpc.NewServer()
	userCtlr := logic.NewUserController(storage.NewSQLDB(db))
	msgCtlr := logic.NewMessageController(storage.NewSQLDB(db))
	api.RegisterChatServer(grpcServer, api.NewChatServer(userCtlr, msgCtlr))
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
	_, err = client.CreateUser(context.Background(), &api.CreateUserRequest{Username: "testuser2", Passphrase: "0123456789abcdef"})
	if err != nil {
		log.Fatalf("Could not create 2nd user account: %v.", err)
	}
	_, err = client.SendMessage(context.Background(),
		&api.SendMessageRequest{Sender: "testuser1", Recipient: "testuser2", Content: "How's it going?"})
	if err != nil {
		log.Fatalf("Could not send a message: %v.", err)
	}
	_, err = client.SendMessage(context.Background(),
		&api.SendMessageRequest{Sender: "testuser2", Recipient: "testuser1", Content: "Can't complain. You?"})
	if err != nil {
		log.Fatalf("Could not send 2nd message: %v.", err)
	}
	_, err = client.SendMessage(context.Background(), &api.SendMessageRequest{
		Sender:    "testuser1",
		Recipient: "testuser2",
		Content:   "https://www.youtube.com/watch?v=9bZkp7q19f0",
	})
	if err != nil {
		log.Fatalf("Could not send 3rd message: %v.", err)
	}
	// Fetch the most recent 2 messages in the conversation.
	conversation, err := client.FetchMessages(context.Background(),
		&api.FetchMessagesRequest{User1: "testuser1", User2: "testuser2", Limit: 2})
	if err != nil {
		log.Fatalf("Could not fetch messages: %v.", err)
	}
	if len(conversation.Messages) == 0 {
		log.Fatalf("No conversation found.")
	}
	if got, want := len(conversation.Messages), 2; got != want {
		log.Fatalf("Conversation has wrong number of messages: got %v, want %v.", got, want)
	}
	if got, want := conversation.Messages[0].Content, "https://www.youtube.com/watch?v=9bZkp7q19f0"; got != want {
		log.Printf("Message content mismatch: got %v, want %v.", got, want)
	}
	if got, want := conversation.Messages[1].Content, "Can't complain. You?"; got != want {
		log.Printf("Message content mismatch: got %v, want %v.", got, want)
	}
	// Now fetch the rest of the conversation.
	conversation, err = client.FetchMessages(context.Background(),
		&api.FetchMessagesRequest{User1: "testuser1", User2: "testuser2", ContinuationToken: conversation.ContinuationToken})
	if err != nil {
		log.Fatalf("Could not fetch messages: %v.", err)
	}
	if len(conversation.Messages) == 0 {
		log.Fatalf("No conversation found.")
	}
	if got, want := len(conversation.Messages), 1; got != want {
		log.Fatalf("Conversation has wrong number of messages: got %v, want %v.", got, want)
	}
	if got, want := conversation.Messages[0].Content, "How's it going?"; got != want {
		log.Printf("Message content mismatch: got %v, want %v.", got, want)
	}
}
