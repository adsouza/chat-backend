package logic

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	AddUser(username string, hash []byte) error
	FetchHash(username string) ([]byte, error)
}

type userController struct {
	db UserStore
}

func NewUserController(db UserStore) *userController {
	return &userController{db: db}
}

func (c *userController) CreateUser(username string, passphrase string) error {
	// Validate that password is long enough.
	if len(passphrase) < 16 {
		return fmt.Errorf("passphrase below 16 char minimum")
	}
	// Check for existing user with identical username.
	if hash, _ := c.db.FetchHash(username); hash != nil {
		//TODO: handle other errors.
		return fmt.Errorf("desired username already taken")
	}
	// Generate a bcrypt hash for the password.
	hash, err := bcrypt.GenerateFromPassword([]byte(passphrase), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("unable to hash passphrase: %v.", err)
	}
	// Persist the username/hash pair to the users table.
	return c.db.AddUser(username, hash)
}

func (c *userController) Authenticate(username string, passphrase string) error {
	hash, err := c.db.FetchHash(username)
	if err != nil {
		return fmt.Errorf("authentication failed because hashed passphrase currently unavailable from storage: %v.", err)
	}
	return bcrypt.CompareHashAndPassword(hash, []byte(passphrase))
}
