package message

import (
	"XAXAtonMTC/pkg/packetsender"
	"XAXAtonMTC/pkg/user"
	"encoding/json"
	"log"
	"time"
)

type Wrapper struct {
	message *Message
}

func NewWrapper(message *Message) *Wrapper {
	return &Wrapper{
		message: message,
	}
}

func (wrapper *Wrapper) NextPacket() (*packetsender.Packet, error) {
	metadata, err := json.Marshal(struct {
		MessageAuthor *user.User `json:"message_author"`
		Created       time.Time  `json:"created"`
	}{
		MessageAuthor: wrapper.message.Author,
		Created:       wrapper.message.Created,
	})
	if err != nil {
		return nil, err
	}

	log.Println("metadata: ", string(metadata))

	if err != nil {
		return nil, err
	}

	return packetsender.NewPacket([]byte(wrapper.message.Text), metadata, packetsender.MESSAGE, false), nil
}
