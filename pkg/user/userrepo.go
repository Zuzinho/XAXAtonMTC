package user

import (
	"github.com/jackc/pgx"
)

type UsersDBRepository struct {
	connConfig pgx.ConnConfig
}

func NewUsersDBRepository(connString string) (*UsersDBRepository, error) {
	conf, err := pgx.ParseConnectionString(connString)
	if err != nil {
		return nil, err
	}

	return &UsersDBRepository{
		connConfig: conf,
	}, nil
}

func (repo *UsersDBRepository) Insert(phoneNumber string) (userID uint32, err error) {
	conn, err := pgx.Connect(repo.connConfig)
	defer conn.Close()
	if err != nil {
		return 0, err
	}

	err = conn.QueryRow("insert into users (phone_number) values ($1) returning user_id;", phoneNumber).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (repo *UsersDBRepository) Update(userID uint32, userName string) error {
	conn, err := pgx.Connect(repo.connConfig)
	defer conn.Close()
	if err != nil {
		return err
	}

	_, err = conn.Exec("update users set user_name = $1 where user_id = $2", userName, userID)

	return err
}

func (repo *UsersDBRepository) SelectByID(userID uint32) (*User, error) {
	conn, err := pgx.Connect(repo.connConfig)
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	var phoneNumber, userName string

	err = conn.QueryRow("select phone_number, user_name from users where user_id = $1", userID).
		Scan(&phoneNumber, &userName)
	if err != nil && phoneNumber == "" {
		return nil, err
	}

	return NewUser(userID, phoneNumber, userName), nil
}

func (repo *UsersDBRepository) SelectByPhoneNumber(phoneNumber string) (*User, error) {
	conn, err := pgx.Connect(repo.connConfig)
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	var userID uint32
	var userName string

	err = conn.QueryRow("select user_id, user_name from users where phone_number = $1", phoneNumber).
		Scan(&userID, &userName)
	if err != nil && userID == 0 {
		return nil, err
	}

	return NewUser(userID, phoneNumber, userName), nil
}
