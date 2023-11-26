package user

import (
	"XAXAtonMTC/pkg/handlers"
	"XAXAtonMTC/pkg/user"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type Handler struct {
	usersRepo user.UsersRepo
}

func NewUserHandler(usersRepo user.UsersRepo) *Handler {
	return &Handler{
		usersRepo: usersRepo,
	}
}

func (handler *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := handlers.UInt32FromMapString(vars, "user_id")
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

	var updateForm updateJSON
	err = json.Unmarshal(body, &updateForm)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	err = handler.usersRepo.Update(userID, updateForm.UserName)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
