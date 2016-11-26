package storage

import (
	"database/sql"
	"fmt"
	"math"
	"time"
)

const (
	PragmaCmd                = "PRAGMA foreign_keys = ON"
	UserTableInitCmd         = "CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY NOT NULL, hash TEXT NOT NULL)"
	ConversationTableInitCmd = `CREATE TABLE IF NOT EXISTS messages (
		timestamp NUMERIC DEFAULT CURRENT_TIMESTAMP NOT NULL, 
		sender TEXT NOT NULL,
		recipient TEXT NOT NULL,
		content TEXT NOT NULL,
		FOREIGN KEY (sender) REFERENCES users(username) ON UPDATE CASCADE ON DELETE RESTRICT,
		FOREIGN KEY (recipient) REFERENCES users(username) ON UPDATE CASCADE ON DELETE RESTRICT)`
)

type Message struct {
	Timestamp       time.Time
	Author, Content string
}

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

func (s *SQLDB) AddMessage(sender, recipient, content string) error {
	_, err := s.Exec("INSERT INTO messages (sender, recipient, content) VALUES (?, ?, ?)", sender, recipient, content)
	return err
}

func (s *SQLDB) ReadMessagesBefore(user1, user2 string, limit uint32, before int64) ([]Message, int64, error) {
	//TODO: use a prepared query.
	rows, err := s.Query(`SELECT rowid, timestamp, sender, content FROM messages WHERE rowid < ? AND sender = ? AND recipient = ?
	UNION ALL SELECT rowid, timestamp, sender, content FROM messages WHERE rowid < ? AND sender = ? AND recipient = ?
	ORDER BY rowid DESC LIMIT ?`,
		before, user1, user2, before, user2, user1, limit)
	if err != nil {
		return nil, math.MaxInt64, fmt.Errorf("unable to execute query for messages between specified users: %v", err)
	}
	defer rows.Close()
	var rowId int64
	var messages []Message
	msg := Message{}
	var ts string
	for rows.Next() {
		err := rows.Scan(&rowId, &ts, &msg.Author, &msg.Content)
		if err != nil {
			return nil, math.MaxInt64, fmt.Errorf("unable to parse data from DB into message struct: %v", err)
		}
		msg.Timestamp, err = time.Parse("2006-01-02 15:04:05", ts)
		if err != nil {
			return nil, math.MaxInt64, fmt.Errorf("unable to parse timestamp from DB: %v", err)
		}
		messages = append(messages, msg)
	}
	return messages, rowId, rows.Err()
}
