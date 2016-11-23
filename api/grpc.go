package api

import (
	"github.com/adsouza/chat-backend/storage"
	"golang.org/x/net/context"
)

type UserController interface {
	CreateUser(username string, passphrase string) error
}

type MessageController interface {
	SendMessage(sender, recipient, message string) error
	FetchMessages(user1, user2 string) ([]storage.Message, error)
}

type chatServer struct {
	userController UserController
	msgController  MessageController
}

func NewChatServer(userCtlr UserController, msgCtlr MessageController) *chatServer {
	return &chatServer{userController: userCtlr, msgController: msgCtlr}
}

func (c *chatServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	return &CreateUserResponse{}, c.userController.CreateUser(req.GetUsername(), req.GetPassphrase())
}

func (c *chatServer) SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
	return &SendMessageResponse{}, c.msgController.SendMessage(req.Sender, req.Recipient, req.Content)
}

func (c *chatServer) FetchMessages(ctx context.Context, req *FetchMessagesRequest) (*FetchMessagesResponse, error) {
	messages, err := c.msgController.FetchMessages(req.User1, req.User2)
	resp := &FetchMessagesResponse{}
	for _, msg := range messages {
		resp.Messages = append(resp.Messages, &Message{Timestamp: msg.Timestamp.Unix(), Author: msg.Author, Content: msg.Content})
	}
	return resp, err
}
