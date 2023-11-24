package handlers

type loginJSON struct {
	PhoneNumber string `json:"phone_number"`
}

type updateJSON struct {
	ID       uint32 `json:"id"`
	UserName string `json:"user_name"`
}
