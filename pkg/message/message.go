package message

import (
	"XAXAtonMTC/pkg/user"
	"time"
)

type Message struct {
	ID      uint32
	Text    string
	Author  *user.User
	Created time.Time
}

func NewMessage(id uint32, text string, author *user.User, created time.Time) *Message {
	return &Message{
		ID:      id,
		Text:    text,
		Author:  author,
		Created: created,
	}
}
