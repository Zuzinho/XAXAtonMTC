package user

type User struct {
	ID          uint32 `json:"id"`
	PhoneNumber string `json:"phone_number"`
	UserName    string `json:"user_name"`
}

func NewUser(id uint32, phoneNumber, userName string) *User {
	return &User{
		ID:          id,
		PhoneNumber: phoneNumber,
		UserName:    userName,
	}
}

type UsersRepo interface {
	Insert(string) (uint32, error)
	Update(uint32, string) error
	SelectByID(uint32) (*User, error)
	SelectByPhoneNumber(string) (*User, error)
}
