package logic

import (
	"github.com/adsouza/chat-backend/storage"
)

type MsgStore interface {
	AddMessage(sender, recipient, content string) error
	ReadMessagesBefore(user1, user2 string, before int64) ([]storage.Message, int64, error)
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

func (c *msgController) FetchMessagesBefore(user1, user2 string, before int64) ([]storage.Message, int64, error) {
	return c.db.ReadMessagesBefore(user1, user2, before)
}
