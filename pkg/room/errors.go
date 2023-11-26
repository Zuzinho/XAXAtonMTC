package room

import "fmt"

type NoRoomError struct {
	id uint32
}

func newNoRoomError(id uint32) NoRoomError {
	return NoRoomError{
		id: id,
	}
}

func (err NoRoomError) Error() string {
	return fmt.Sprintf("no room by id %d", err.id)
}

var NoRoomErr NoRoomError = NoRoomError{}
