package logic

import (
	"fmt"

	"github.com/adsouza/chat-backend/storage"
)

const DELIM = ':'

type MsgStore interface {
	AddMessage(conversationId, author, content string) error
	ReadMessages(conversationId string) ([]storage.Message, error)
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

func conversationIdFromParticipants(user1, user2 string) string {
	if user1 < user2 {
		return fmt.Sprintf("%s%c%s", user1, DELIM, user2)
	}
	return fmt.Sprintf("%s%c%s", user2, DELIM, user1)
}

func (c *msgController) SendMessage(sender, recipient, message string) error {
	return c.db.AddMessage(conversationIdFromParticipants(sender, recipient), sender, message)
}

func (c *msgController) FetchMessages(user1, user2 string) ([]storage.Message, error) {
	return c.db.ReadMessages(conversationIdFromParticipants(user1, user2))
}
