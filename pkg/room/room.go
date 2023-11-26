package room

import (
	"XAXAtonMTC/pkg/message"
	"XAXAtonMTC/pkg/user"
)

type AccessType int

const (
	PRIVATE AccessType = 0
	PUBLIC  AccessType = 1
)

type Room struct {
	ID       uint32                `json:"id"`
	Token    string                `json:"-"`
	Name     string                `json:"name"`
	Author   *user.User            `json:"author"`
	Users    map[uint32]*user.User `json:"users"`
	Messages []*message.Message    `json:"messages"`
	Access   AccessType            `json:"access"`
}

func NewRoom(roomID uint32, token string, name string, author *user.User, accessType AccessType) *Room {
	return &Room{
		ID:       roomID,
		Token:    token,
		Name:     name,
		Author:   author,
		Users:    make(map[uint32]*user.User),
		Messages: make([]*message.Message, 0),
		Access:   accessType,
	}
}

type RoomsRepo interface {
	CreateRoom(*user.User, string, string, AccessType) *Room
	AddListener(uint32, *user.User) (*Room, error)
	DeleteListener(uint32, uint32) (*Room, error)
	UsersID(uint32) ([]uint32, error)
	AccessedRooms() []*Room
	DeleteRoom(uint32)
	AddMessage(uint32, string, *user.User) (*message.Message, error)
	Access(uint32) (AccessType, error)
	IsValidToken(uint32, string) (bool, error)
}
