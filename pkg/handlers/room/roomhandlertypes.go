package room

import "XAXAtonMTC/pkg/room"

type roomPostType struct {
	UserID uint32          `json:"user_id"`
	Token  string          `json:"token"`
	Name   string          `json:"name"`
	Access room.AccessType `json:"access"`
}

type roomGetType struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	UsersCount int    `json:"users_count"`
}
