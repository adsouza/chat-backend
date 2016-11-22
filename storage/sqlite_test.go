package storage_test

import (
	"database/sql"
	"testing"

	"github.com/adsouza/chat-backend/storage"
	_ "github.com/mattn/go-sqlite3"
)

func newStore(t *testing.T) (*storage.SQLDB, func()) {
	db, err := sql.Open("sqlite3", "")
	if err != nil {
		t.Fatalf("Unable to open connection to DB: %v.", err)
	}
	if _, err := db.Exec(storage.UserTableInitCmd); err != nil {
		db.Close()
		t.Fatalf("Unable to create new users table in test DB: %v.", err)
	}
	return storage.NewSQLDB(db), func() { db.Close() }
}

func TestHappyPath(t *testing.T) {
	store, closer := newStore(t)
	defer closer()
	if err := store.AddUser("testuser1", []byte("012345678901234567890123456789012345678901234567890123456789")); err != nil {
		t.Errorf("Unable to add a new row to the users table: %v.", err)
	}
	hash, err := store.FetchHash("testuser1")
	if err != nil {
		t.Errorf("Unable to retrieve hash for specified user: %v.", err)
	}
	if got, want := string(hash), "012345678901234567890123456789012345678901234567890123456789"; got != want {
		t.Errorf("Hash mismatch:\ngot  %v\nwant %v", got, want)
	}
}

func TestNonexistentUser(t *testing.T) {
	store, closer := newStore(t)
	defer closer()
	if _, err := store.FetchHash("testuser1"); err == nil {
		t.Errorf("Able to retrieve hash for nonexistent user!")
	}
}
