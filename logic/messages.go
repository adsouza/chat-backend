package logic

import (
	"net/url"
	"strings"

	"github.com/adsouza/chat-backend/storage"
)

const (
	Vevo = "www.vevo.com/watch"
)

type MsgStore interface {
	AddMessage(sender, recipient, content string, metadata []byte) error
	ReadMessagesBefore(user1, user2 string, limit uint32, before int64) ([]storage.Message, int64, error)
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
	url, err := url.Parse(message)
	if err == nil {
		if (url.Host == "www.youtube.com" || url.Host == "www.vevo.com") && strings.HasPrefix(url.Path, "/watch/") {
			// We have a video!
		}
		// Assume we have an image.
	}
	return c.db.AddMessage(sender, recipient, message, nil)
}

func (c *msgController) FetchMessagesBefore(user1, user2 string, limit uint32, before int64) ([]storage.Message, int64, error) {
	return c.db.ReadMessagesBefore(user1, user2, limit, before)
}
