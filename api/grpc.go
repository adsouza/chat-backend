package api

import (
	"fmt"
	"math"

	"github.com/adsouza/chat-backend/storage"
	"golang.org/x/net/context"
)

type UserController interface {
	CreateUser(username string, passphrase string) error
}

type MessageController interface {
	SendMessage(sender, recipient, message string) error
	FetchMessagesBefore(user1, user2 string, limit uint32, before int64) ([]storage.Message, int64, error)
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
	if req.User1 == "" || req.User2 == "" {
		return &FetchMessagesResponse{}, fmt.Errorf("both the User1 & User2 fields are required")
	}
	before := req.ContinuationToken
	if before == 0 {
		before = math.MaxInt64
	}
	limit := req.Limit
	if limit == 0 {
		limit = math.MaxUint32
	}
	messages, continuationToken, err := c.msgController.FetchMessagesBefore(req.User1, req.User2, limit, before)
	resp := &FetchMessagesResponse{ContinuationToken: continuationToken}
	for _, msg := range messages {
		resp.Messages = append(resp.Messages, &Message{Timestamp: msg.Timestamp.Unix(), Author: msg.Author, Content: msg.Content})
	}
	return resp, err
}
