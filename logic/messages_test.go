package logic_test

import (
	"fmt"
	"testing"

	"github.com/adsouza/chat-backend/logic"
	"github.com/adsouza/chat-backend/storage"
)

type mockMsgStore struct {
	conversations map[string][]storage.Message
}

func (m *mockMsgStore) AddMessage(conversationId, author, content string) error {
	m.conversations[conversationId] = append(m.conversations[conversationId], storage.Message{Author: author, Content: content})
	return nil
}

func (m *mockMsgStore) ReadMessages(conversationId string) ([]storage.Message, error) {
	conversation, ok := m.conversations[conversationId]
	if !ok {
		return nil, fmt.Errorf("no row with key %v exists", conversationId)
	}
	return conversation, nil
}

type mockDb struct {
	mockUserStore
	mockMsgStore
}

func TestHappyPath(t *testing.T) {
	mockDb := &mockDb{
		mockUserStore: mockUserStore{hashes: make(map[string][]byte)},
		mockMsgStore:  mockMsgStore{conversations: make(map[string][]storage.Message)},
	}
	userCtlr := logic.NewUserController(&mockDb.mockUserStore)
	if err := userCtlr.CreateUser("testuser1", "123456789abcdefg"); err != nil {
		t.Fatalf("16 char passphrase was not permitted but should be.")
	}
	if err := userCtlr.CreateUser("testuser2", "123456789abcdefg"); err != nil {
		t.Errorf("2nd user account was not permitted but should be.")
	}
	msgCtlr := logic.NewMessageController(mockDb)
	if err := msgCtlr.SendMessage("testuser1", "testuser2", "Bonjour!"); err != nil {
		t.Fatalf("Sending a message failed: %v.", err)
	}
	conversation, err := msgCtlr.FetchMessages("testuser1", "testuser2")
	if err != nil {
		t.Fatalf("Unable to fetch a conversation: %v.", err)
	}
	if len(conversation) < 1 {
		t.Fatalf("No conversation found.")
	}
}
