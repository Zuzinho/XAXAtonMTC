package message

import (
	"XAXAtonMTC/pkg/handlers"
	"XAXAtonMTC/pkg/message"
	"XAXAtonMTC/pkg/room"
	"XAXAtonMTC/pkg/sender"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type Handler struct {
	RoomsRepo   room.RoomsRepo
	SendersRepo sender.SendersRepo
}

func NewHandler(roomsRepo room.RoomsRepo, sendersRepo sender.SendersRepo) *Handler {
	return &Handler{
		RoomsRepo:   roomsRepo,
		SendersRepo: sendersRepo,
	}
}

func (handler *Handler) Message(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roomID, err := handlers.UInt32FromMapString(vars, "room_id")
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	var messagePostForm messagePostType

	err = json.Unmarshal(body, &messagePostForm)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	msg, err := handler.RoomsRepo.AddMessage(roomID, messagePostForm.Text, &messagePostForm.Author)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	wrapper := message.NewWrapper(msg)

	userIDs, err := handler.RoomsRepo.UsersID(roomID)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	packet, err := wrapper.NextPacket()
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	handler.SendersRepo.SendPacket(userIDs, func(userID uint32) bool {
		return userID != messagePostForm.Author.ID
	}, packet)

	w.WriteHeader(http.StatusOK)
}
