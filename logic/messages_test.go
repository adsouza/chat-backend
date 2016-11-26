package logic_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/adsouza/chat-backend/logic"
	"github.com/adsouza/chat-backend/storage"
)

func conversationIdFromParticipants(user1, user2 string) string {
	if user1 < user2 {
		return fmt.Sprintf("%s%c%s", user1, ':', user2)
	}
	return fmt.Sprintf("%s%c%s", user2, ':', user1)
}

type mockMsgStore struct {
	conversations map[string][]storage.Message
}

func (m *mockMsgStore) AddMessage(sender, recipient, content string) error {
	// Add the new message to the beginning.
	conversationId := conversationIdFromParticipants(sender, recipient)
	m.conversations[conversationId] = append([]storage.Message{storage.Message{Author: sender, Content: content}}, m.conversations[conversationId]...)
	return nil
}

func (m *mockMsgStore) ReadMessagesBefore(user1, user2 string, before int64) ([]storage.Message, int64, error) {
	conversationId := conversationIdFromParticipants(user1, user2)
	conversation, ok := m.conversations[conversationId]
	if !ok {
		return nil, math.MaxInt64, fmt.Errorf("no row with key %v exists", conversationId)
	}
	return conversation, math.MaxInt64, nil
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
	conversation, _, err := msgCtlr.FetchMessagesBefore("testuser1", "testuser2", math.MaxInt64)
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
	conversation, _, err = msgCtlr.FetchMessagesBefore("testuser2", "testuser1", math.MaxInt64)
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
