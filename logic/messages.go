package logic

import (
	"github.com/adsouza/chat-backend/storage"
)

type MsgStore interface {
	AddMessage(sender, recipient, content string) error
	ReadMessages(user1, user2 string) ([]storage.Message, error)
}

type Db interface {
	UserStore
	MsgStore
}

type msgController struct {
	db Db
}

func NewMessageController(db Db) *msgController {
	return &msgController{db: db}
}

func (c *msgController) SendMessage(sender, recipient, message string) error {
	return c.db.AddMessage(sender, recipient, message)
}

func (c *msgController) FetchMessages(user1, user2 string) ([]storage.Message, error) {
	return c.db.ReadMessages(user1, user2)
}
