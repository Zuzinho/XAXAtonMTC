package handlers

import (
	"XAXAtonMTC/pkg/user"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	"io"
	"log"
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

func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpJSONErr(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var loginForm loginJSON
	err = json.Unmarshal(body, &loginForm)
	if err != nil {
		httpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	log.Println(loginForm.PhoneNumber)

	var u *user.User

	u, err = handler.usersRepo.SelectByPhoneNumber(loginForm.PhoneNumber)
	if errors.As(err, &pgx.ErrNoRows) {
		userID, err := handler.usersRepo.Insert(loginForm.PhoneNumber)
		if err != nil {
			httpJSONErr(w, err, http.StatusInternalServerError)
			return
		}

		u = user.NewUser(userID, loginForm.PhoneNumber, "")
	} else if err != nil {
		httpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	buf, err := json.Marshal(*u)
	if err != nil {
		httpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		httpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

func httpJSONErr(w http.ResponseWriter, err error, status int) {
	http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err.Error()), status)
}
