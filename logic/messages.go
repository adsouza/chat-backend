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

func metadataFromURL(url *url.URL) *api.Metadata {
	if strings.HasPrefix(url.Path, "/watch/") {
		switch url.Host {
		case "www.youtube.com":
			return &api.Metadata{Media: &api.Metadata_Video{Video: &api.Video{Source: api.Video_YOUTUBE}}}
		case "www.vevo.com":
			return &api.Metadata{Media: &api.Metadata_Video{Video: &api.Video{Source: api.Video_VEVO}}}
		default:
		}
	}
	// Assume we have an image.
	return &api.Metadata{Media: &api.Metadata_Image{Image: &api.Image{}}}
}

func (c *msgController) SendMessage(sender, recipient, message string) error {
	var data []byte
	if url, err := url.Parse(message); err == nil {
		data, err = proto.Marshal(metadataFromURL(url))
		if err != nil {
			return fmt.Errorf("could not marshal metadata proto into blob: %v", err)
		}
	}
	return c.db.AddMessage(sender, recipient, message, data)
}

func (c *msgController) FetchMessagesBefore(user1, user2 string, limit uint32, before int64) ([]storage.Message, int64, error) {
	return c.db.ReadMessagesBefore(user1, user2, limit, before)
}
