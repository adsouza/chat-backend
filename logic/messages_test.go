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
	// Add the new message to the beginning.
	m.conversations[conversationId] = append([]storage.Message{storage.Message{Author: author, Content: content}}, m.conversations[conversationId]...)
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
		t.Fatalf("2nd user account was not permitted but should be.")
	}
	msgCtlr := logic.NewMessageController(mockDb)
	if err := msgCtlr.SendMessage("testuser1", "testuser2", "Bonjour!"); err != nil {
		t.Fatalf("Sending a message failed: %v.", err)
	}
	if err := msgCtlr.SendMessage("testuser2", "testuser1", "A revoir."); err != nil {
		t.Fatalf("Sending a 2nd message failed: %v.", err)
	}
	conversation, err := msgCtlr.FetchMessages("testuser1", "testuser2")
	if err != nil {
		t.Fatalf("Unable to fetch a conversation: %v.", err)
	}
	if conversation == nil || len(conversation) == 0 {
		t.Fatalf("No conversation found.")
	}
	if got, want := len(conversation), 2; got != want {
		t.Fatalf("Conversation has wrong number of messages: got %v, want %v.", got, want)
	}
	if got, want := conversation[0].Content, "A revoir."; got != want {
		t.Errorf("Message content mismatch: got %v, want %v.", got, want)
	}
	if got, want := conversation[1].Content, "Bonjour!"; got != want {
		t.Errorf("Message content mismatch: got %v, want %v.", got, want)
	}
	// Now make sure it works with the usernames in reverse order too.
	conversation, err = msgCtlr.FetchMessages("testuser2", "testuser1")
	if err != nil {
		t.Fatalf("Unable to fetch a conversation: %v.", err)
	}
	if conversation == nil || len(conversation) == 0 {
		t.Fatalf("No conversation found.")
	}
	if got, want := len(conversation), 2; got != want {
		t.Fatalf("Conversation has wrong number of messages: got %v, want %v.", got, want)
	}
	if got, want := conversation[0].Content, "A revoir."; got != want {
		t.Errorf("Message content mismatch: got %v, want %v.", got, want)
	}
	if got, want := conversation[1].Content, "Bonjour!"; got != want {
		t.Errorf("Message content mismatch: got %v, want %v.", got, want)
	}
}
