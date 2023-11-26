package connection

import (
	"XAXAtonMTC/pkg/handlers"
	"XAXAtonMTC/pkg/room"
	"XAXAtonMTC/pkg/sender"
	"github.com/gorilla/mux"
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

func (handler *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := handlers.UInt32FromMapString(vars, "user_id")
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	err = handler.SendersRepo.AddSender(userID, w, r)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusSwitchingProtocols)
}
