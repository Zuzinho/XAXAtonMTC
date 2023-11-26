package room

import (
	"XAXAtonMTC/pkg/handlers"
	"XAXAtonMTC/pkg/room"
	"XAXAtonMTC/pkg/sender"
	"XAXAtonMTC/pkg/user"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	UsersRepo   user.UsersRepo
	RoomsRepo   room.RoomsRepo
	SendersRepo sender.SendersRepo
}

func NewHandler(usersRepo user.UsersRepo, roomsRepo room.RoomsRepo, sendersRepo sender.SendersRepo) *Handler {
	return &Handler{
		UsersRepo:   usersRepo,
		RoomsRepo:   roomsRepo,
		SendersRepo: sendersRepo,
	}
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var roomPostForm roomPostType

	err = json.Unmarshal(body, &roomPostForm)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	author, err := handler.UsersRepo.SelectByID(roomPostForm.UserID)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	rm := handler.RoomsRepo.CreateRoom(author, roomPostForm.Token, roomPostForm.Name, roomPostForm.Access)

	buf, err := json.Marshal(*rm)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) AddListener(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roomID, err := handlers.UInt32FromMapString(vars, "room_id")
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	access, err := handler.RoomsRepo.Access(roomID)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	if access == room.PRIVATE {
		handlers.HttpJSONErr(w, err, http.StatusUnauthorized)
		return
	}

	userID, err := handlers.UInt32FromMapString(vars, "user_id")
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	u, err := handler.UsersRepo.SelectByID(userID)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	rm, err := handler.RoomsRepo.AddListener(roomID, u)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	buf, err := json.Marshal(*rm)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) AddListenerByRef(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roomID, err := handlers.UInt32FromMapString(vars, "room_id")
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	token := vars["token"]

	access, err := handler.RoomsRepo.Access(roomID)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	if access == room.PRIVATE {
		valid, err := handler.RoomsRepo.IsValidToken(roomID, token)
		if err != nil {
			handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
			return
		}

		if !valid {
			handlers.HttpJSONErr(w, err, http.StatusUnauthorized)
			return
		}
	}

	userID, err := handlers.UInt32FromMapString(vars, "user_id")
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	u, err := handler.UsersRepo.SelectByID(userID)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	rm, err := handler.RoomsRepo.AddListener(roomID, u)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	buf, err := json.Marshal(*rm)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) DeleteListener(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roomID, err := handlers.UInt32FromMapString(vars, "room_id")
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	userID, err := handlers.UInt32FromMapString(vars, "user_id")
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	rm, err := handler.RoomsRepo.DeleteListener(roomID, userID)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	handler.SendersRepo.DeleteSender(userID)

	if rm.Author.ID == userID {
		go func(roomID uint32) {
			usersID, err := handler.RoomsRepo.UsersID(roomID)
			if err != nil {
				log.Println(err)
				return
			}
			go handler.SendersRepo.CloseConnections(usersID...)
		}(roomID)
		go handler.RoomsRepo.DeleteRoom(rm.ID)
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) Rooms(w http.ResponseWriter, r *http.Request) {
	rooms := handler.RoomsRepo.AccessedRooms()

	roomForms := make([]roomGetType, len(rooms))
	for i, rm := range rooms {
		roomForms[i] = roomGetType{
			ID:         rm.ID,
			Name:       rm.Name,
			UsersCount: len(rm.Users) + 1,
		}
	}

	buf, err := json.Marshal(roomForms)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
