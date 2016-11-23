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
