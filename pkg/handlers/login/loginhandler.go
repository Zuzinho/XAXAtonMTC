package login

import (
	"XAXAtonMTC/pkg/handlers"
	"XAXAtonMTC/pkg/sms"
	"XAXAtonMTC/pkg/user"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	usersRepo   user.UsersRepo
	smsSender   *sms.Sender
	phoneNumber string
	smsCode     string
}

func NewHandler(usersRepo user.UsersRepo, smsApiID string) *Handler {
	return &Handler{
		usersRepo: usersRepo,
		smsSender: sms.NewSender(smsApiID),
	}
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var loginForm loginJSON
	err = json.Unmarshal(body, &loginForm)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	log.Println(loginForm.PhoneNumber)

	handler.phoneNumber = loginForm.PhoneNumber

	handler.smsCode = generateRandomCode()
	err = handler.smsSender.SendSms(handler.phoneNumber, handler.smsCode)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) CheckSms(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var smsCodeForm smsCodeJSON
	err = json.Unmarshal(body, &smsCodeForm)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	if smsCodeForm.SmsCode != handler.smsCode {
		handlers.HttpJSONErr(w, InvalidSmsCodeErr, http.StatusBadRequest)
		return
	}

	u, err := handler.usersRepo.SelectByPhoneNumber(handler.phoneNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		userID, err := handler.usersRepo.Insert(handler.phoneNumber)
		if err != nil {
			handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
			return
		}

		u = user.NewUser(userID, handler.phoneNumber, "")
	} else if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
		return
	}

	buf, err := json.Marshal(*u)
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
