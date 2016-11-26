package storage_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/adsouza/chat-backend/storage"
	_ "github.com/mattn/go-sqlite3"
)

func newStore(t *testing.T) (*storage.SQLDB, func()) {
	db, err := sql.Open("sqlite3", "")
	if err != nil {
		t.Fatalf("Unable to open connection to DB: %v.", err)
	}
	if _, err := db.Exec(storage.PragmaCmd); err != nil {
		t.Errorf("Unable to enable foreign key constraints in test DB: %v.", err)
	}
	if _, err := db.Exec(storage.UserTableInitCmd); err != nil {
		db.Close()
		t.Fatalf("Unable to create new users table in test DB: %v.", err)
	}
	if _, err := db.Exec(storage.ConversationTableInitCmd); err != nil {
		db.Close()
		t.Fatalf("Unable to create new conversations table in test DB: %v.", err)
	}
	return storage.NewSQLDB(db), func() { db.Close() }
}

func TestHappyPath(t *testing.T) {
	store, closer := newStore(t)
	defer closer()
	if err := store.AddUser("testuser1", []byte("012345678901234567890123456789012345678901234567890123456789")); err != nil {
		t.Fatalf("Unable to add a new row to the users table: %v.", err)
	}
	hash, err := store.FetchHash("testuser1")
	if err != nil {
		t.Fatalf("Unable to retrieve hash for recently added user: %v.", err)
	}
	if got, want := string(hash), "012345678901234567890123456789012345678901234567890123456789"; got != want {
		t.Errorf("Hash mismatch:\ngot  %v\nwant %v", got, want)
	}
	if err := store.AddUser("testuser2", []byte("012345678901234567890123456789012345678901234567890123456789")); err != nil {
		t.Fatalf("Unable to add a 2nd row to the users table: %v.", err)
	}
	if err := store.AddMessage("testuser1", "testuser2", "Hello!"); err != nil {
		t.Fatalf("Unable to add a new row to the messages table: %v.", err)
	}
	messages, err := store.ReadMessages("testuser1", "testuser2")
	if err != nil {
		t.Fatalf("Unable to retrieve messages for specified conversation: %v.", err)
	}
	if messages == nil {
		t.Fatalf("No messages found for recently initiated conversation.")
	}
	if got, want := messages[0].Content, "Hello!"; got != want {
		t.Errorf("Message content mismatch: got %v, want %v.", got, want)
	}
	if got, want := messages[0].Author, "testuser1"; got != want {
		t.Errorf("Message author mismatch: got %v, want %v.", got, want)
	}
}

func TestNonexistentUser(t *testing.T) {
	store, closer := newStore(t)
	defer closer()
	if _, err := store.FetchHash("testuser1"); err == nil {
		t.Errorf("Able to retrieve hash for nonexistent user!")
	}
}

func TestMessageOrder(t *testing.T) {
	store, closer := newStore(t)
	defer closer()
	if err := store.AddUser("testuser1", []byte("012345678901234567890123456789012345678901234567890123456789")); err != nil {
		t.Fatalf("Unable to add a new row to the users table: %v.", err)
	}
	if err := store.AddUser("testuser2", []byte("012345678901234567890123456789012345678901234567890123456789")); err != nil {
		t.Fatalf("Unable to add a new row to the users table: %v.", err)
	}
	if err := store.AddMessage("testuser1", "testuser2", "Hello!"); err != nil {
		t.Fatalf("Unable to add a new row to the messages table: %v.", err)
	}
	time.Sleep(time.Second)
	if err := store.AddMessage("testuser2", "testuser1", "Goodbye."); err != nil {
		t.Fatalf("Unable to add a 2nd row to the messages table: %v.", err)
	}
	messages, err := store.ReadMessages("testuser1", "testuser2")
	if err != nil {
		t.Fatalf("Unable to retrieve messages for specified conversation: %v.", err)
	}
	if messages == nil {
		t.Fatalf("No messages found for recently initiated conversation.")
	}
	if got, want := len(messages), 2; got != want {
		t.Fatalf("Wrong number of messages retrieved: got %v, want %v.", got, want)
	}
	if got, want := messages[0].Content, "Goodbye."; got != want {
		t.Errorf("Message content mismatch: got %v, want %v.", got, want)
	}
	if got, want := messages[1].Content, "Hello!"; got != want {
		t.Errorf("Message content mismatch: got %v, want %v.", got, want)
	}
}

func TestMsgFromNonexistentUser(t *testing.T) {
	store, closer := newStore(t)
	defer closer()
	if err := store.AddUser("testuser1", []byte("012345678901234567890123456789012345678901234567890123456789")); err != nil {
		t.Fatalf("Unable to add a new row to the users table: %v.", err)
	}
	if err := store.AddMessage("testuser2", "testuser1", "Hello!"); err == nil {
		t.Errorf("Able to add a new row to the messages table with a nonexistent sender!")
	}
}

func TestMsgToNonexistentUser(t *testing.T) {
	store, closer := newStore(t)
	defer closer()
	if err := store.AddUser("testuser1", []byte("012345678901234567890123456789012345678901234567890123456789")); err != nil {
		t.Fatalf("Unable to add a new row to the users table: %v.", err)
	}
	if err := store.AddMessage("testuser1", "testuser2", "Hello!"); err == nil {
		t.Errorf("Able to add a new row to the messages table with a nonexistent recipient!")
	}
}
