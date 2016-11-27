package logic

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/adsouza/chat-backend/api"
	"github.com/adsouza/chat-backend/storage"
	"github.com/golang/protobuf/proto"
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
	metadata := &api.Metadata{}
	url, err := url.Parse(message)
	if err == nil {
		if (url.Host == "www.youtube.com" || url.Host == "www.vevo.com") && strings.HasPrefix(url.Path, "/watch/") {
			// We have a video!
		}
		// Assume we have an image.
	}
	data, err := proto.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("could not marshal metadata proto into blob: %v", err)
	}
	return c.db.AddMessage(sender, recipient, message, data)
}

func (c *msgController) FetchMessagesBefore(user1, user2 string, limit uint32, before int64) ([]storage.Message, int64, error) {
	return c.db.ReadMessagesBefore(user1, user2, limit, before)
}
