package storage

import (
	"database/sql"
	"fmt"
)

const UserTableInitCmd = "CREATE TABLE users (username TEXT PRIMARY KEY NOT NULL, hash TEXT NOT NULL)"

type SQLDB struct {
	*sql.DB
}

func NewSQLDB(db *sql.DB) *SQLDB {
	return &SQLDB{db}
}

func (s *SQLDB) AddUser(username string, hash []byte) error {
	_, err := s.Exec("INSERT INTO users (username, hash) VALUES (?, ?)", username, hash)
	return err
}

func (s *SQLDB) FetchHash(username string) ([]byte, error) {
	var hash string
	err := s.QueryRow("SELECT hash FROM users WHERE username=?", username).Scan(&hash)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("no such username found")
	case err != nil:
		return nil, fmt.Errorf("unexpected DB access failure: %v", err)
	default:
		return []byte(hash), nil
	}
}
