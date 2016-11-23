package storage

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	UserTableInitCmd         = "CREATE TABLE users (username TEXT PRIMARY KEY NOT NULL, hash TEXT NOT NULL)"
	ConversationTableInitCmd = "CREATE TABLE messages (conversationid TEXT NOT NULL, timestamp NUMERIC DEFAULT CURRENT_TIMESTAMP NOT NULL, author TEXT NOT NULL, content TEXT NOT NULL)"
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

func (s *SQLDB) AddMessage(conversationId, author, content string) error {
	_, err := s.Exec("INSERT INTO messages (conversationid, author, content) VALUES (?, ?, ?)", conversationId, author, content)
	return err
}

func (s *SQLDB) ReadMessages(conversationId string) ([]Message, error) {
	//TODO: use a prepared query.
	rows, err := s.Query("SELECT timestamp, author, content FROM messages WHERE conversationid = ? ORDER BY timestamp DESC", conversationId)
	if err != nil {
		return nil, fmt.Errorf("unable to execute query for messages in specified conversation: %v", err)
	}
	defer rows.Close()
	var messages []Message
	msg := Message{}
	var ts string
	for rows.Next() {
		err := rows.Scan(&ts, &msg.Author, &msg.Content)
		if err != nil {
			return nil, fmt.Errorf("unable to parse data from DB into message struct: %v", err)
		}
		msg.Timestamp, err = time.Parse("2006-01-02 15:04:05", ts)
		if err != nil {
			return nil, fmt.Errorf("unable to parse timestamp from DB: %v", err)
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}
