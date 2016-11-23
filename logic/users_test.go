package logic_test

import (
	"fmt"
	"testing"

	"github.com/adsouza/chat-backend/logic"
)

type mockUserStore struct {
	hashes map[string][]byte
}

func (m *mockUserStore) AddUser(username string, hash []byte) error {
	m.hashes[username] = hash
	return nil
}

func (m *mockUserStore) FetchHash(username string) ([]byte, error) {
	hash, ok := m.hashes[username]
	if !ok {
		return nil, fmt.Errorf("no row with key %v exists", username)
	}
	return hash, nil
}

func TestUsersHappyPath(t *testing.T) {
	userCtlr := logic.NewUserController(&mockUserStore{hashes: make(map[string][]byte)})
	if err := userCtlr.CreateUser("testuser1", "123456789abcdefg"); err != nil {
		t.Fatalf("16 char passphrase was not permitted but should be.")
	}
	if err := userCtlr.CreateUser("testuser2", "123456789abcdefg"); err != nil {
		t.Errorf("2nd user account was not permitted but should be.")
	}
	if err := userCtlr.Authenticate("testuser1", "123456789abcdefg"); err != nil {
		t.Errorf("Unable to authenticate user that was just added.")
	}
}

func TestShortPassphrase(t *testing.T) {
	userCtlr := logic.NewUserController(&mockUserStore{})
	if err := userCtlr.CreateUser("testuser1", "123456789abcdef"); err == nil {
		t.Errorf("Passphrase shorter than 16 chars was permitted but should not be.")
	}
}

func TestDupeUsername(t *testing.T) {
	userCtlr := logic.NewUserController(&mockUserStore{hashes: make(map[string][]byte)})
	if err := userCtlr.CreateUser("testuser1", "123456789abcdefg"); err != nil {
		t.Errorf("16 char passphrase was not permitted but should be.")
	}
	if err := userCtlr.CreateUser("testuser1", "123456789abcdefg"); err == nil {
		t.Errorf("Duplicate username was permitted but should not be.")
	}
}

func TestNonexistentUser(t *testing.T) {
	userCtlr := logic.NewUserController(&mockUserStore{})
	if err := userCtlr.Authenticate("testuser1", "123456789abcdefg"); err == nil {
		t.Errorf("Managed to authenticate user that was never added!")
	}
}

func TestWrongPassphrase(t *testing.T) {
	userCtlr := logic.NewUserController(&mockUserStore{hashes: make(map[string][]byte)})
	if err := userCtlr.CreateUser("testuser1", "123456789abcdefg"); err != nil {
		t.Errorf("16 char passphrase was not permitted but should be.")
	}
	if err := userCtlr.Authenticate("testuser1", "123456789abcdef!"); err == nil {
		t.Errorf("Managed to authenticate user using wrong passphrase!")
	}
}
