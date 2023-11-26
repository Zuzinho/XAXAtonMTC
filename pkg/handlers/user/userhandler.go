package handlers

import (
	"XAXAtonMTC/pkg/user"
	"encoding/json"
	"io"
	"net/http"
)

type UserHandler struct {
	usersRepo user.UsersRepo
}

func NewUserHandler(usersRepo user.UsersRepo) *UserHandler {
	return &UserHandler{
		usersRepo: usersRepo,
	}
}

func (handler *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	var updateForm updateJSON
	err = json.Unmarshal(body, &updateForm)
	if err != nil {
		httpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	err = handler.usersRepo.Update(updateForm.ID, updateForm.UserName)
	if err != nil {
		httpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
