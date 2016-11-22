package api

import (
	"golang.org/x/net/context"
)

type UserController interface {
	CreateUser(username string, passphrase string) error
}

type chatServer struct {
	controller UserController
}

func NewChatServer(ctlr UserController) *chatServer {
	return &chatServer{controller: ctlr}
}

func (c *chatServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	return &CreateUserResponse{}, c.controller.CreateUser(req.GetUsername(), req.GetPassphrase())
}
