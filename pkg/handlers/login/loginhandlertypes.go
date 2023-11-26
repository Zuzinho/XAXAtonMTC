package login

type loginJSON struct {
	PhoneNumber string `json:"phone_number"`
}

type smsCodeJSON struct {
	SmsCode string `json:"code"`
}
