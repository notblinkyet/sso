package models

type User struct {
	Id       int64
	Email    string
	PassHash []byte
}

func NewUser(id int64, email string, passHash []byte) *User {
	return &User{
		Id:       id,
		Email:    email,
		PassHash: passHash,
	}
}
