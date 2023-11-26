package message

import "XAXAtonMTC/pkg/user"

type messagePostType struct {
	Author user.User `json:"author"`
	Text   string    `json:"text"`
}
