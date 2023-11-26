package room

import (
	"XAXAtonMTC/pkg/message"
	"XAXAtonMTC/pkg/user"
	"sync"
	"time"
)

type RoomsMemoryRepository struct {
	mu    *sync.RWMutex
	rooms map[uint32]*Room
	newID uint32
}

func NewRoomsMemoryRepository() *RoomsMemoryRepository {
	return &RoomsMemoryRepository{
		mu:    &sync.RWMutex{},
		rooms: make(map[uint32]*Room),
		newID: 1,
	}
}

func (repo *RoomsMemoryRepository) CreateRoom(author *user.User, token, name string, accessType AccessType) *Room {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	room := NewRoom(repo.newID, token, name, author, accessType)

	repo.rooms[repo.newID] = room

	repo.newID++

	return room
}

func (repo *RoomsMemoryRepository) AddListener(roomID uint32, newListener *user.User) (*Room, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	room, ok := repo.rooms[roomID]
	if !ok {
		return nil, newNoRoomError(roomID)
	}

	room.Users[newListener.ID] = newListener

	return room, nil
}

func (repo *RoomsMemoryRepository) DeleteListener(roomID uint32, oldListenerID uint32) (*Room, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	room, ok := repo.rooms[roomID]
	if !ok {
		return nil, newNoRoomError(roomID)
	}

	delete(room.Users, oldListenerID)

	return room, nil
}

func (repo *RoomsMemoryRepository) UsersID(roomID uint32) ([]uint32, error) {
	repo.mu.RLock()
	room, ok := repo.rooms[roomID]
	repo.mu.RUnlock()

	if !ok {
		return nil, newNoRoomError(roomID)
	}

	userIDs := make([]uint32, len(room.Users)+1)

	userIDs[0] = room.Author.ID

	i := 1
	for k := range room.Users {
		userIDs[i] = k
		i++
	}

	return userIDs, nil
}

func (repo *RoomsMemoryRepository) AccessedRooms() []*Room {
	rooms := make([]*Room, 0)

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, v := range repo.rooms {
		if v.Access == PUBLIC {
			rooms = append(rooms, v)
		}
	}

	return rooms
}

func (repo *RoomsMemoryRepository) DeleteRoom(roomID uint32) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	delete(repo.rooms, roomID)
}

func (repo *RoomsMemoryRepository) AddMessage(roomID uint32, text string, author *user.User) (*message.Message, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	rm, ok := repo.rooms[roomID]

	if !ok {
		return nil, NoRoomErr
	}

	msg := message.NewMessage(uint32(len(rm.Messages)+1), text, author, time.Now())

	rm.Messages = append(rm.Messages, msg)

	return msg, nil
}

func (repo *RoomsMemoryRepository) Access(roomID uint32) (AccessType, error) {
	repo.mu.RLock()
	rm, ok := repo.rooms[roomID]
	repo.mu.RUnlock()

	if !ok {
		return 0, NoRoomErr
	}

	return rm.Access, nil
}

func (repo *RoomsMemoryRepository) IsValidToken(roomID uint32, token string) (bool, error) {
	repo.mu.RLock()
	rm, ok := repo.rooms[roomID]
	repo.mu.RUnlock()

	if !ok {
		return false, NoRoomErr
	}

	return rm.Token == token, nil
}
