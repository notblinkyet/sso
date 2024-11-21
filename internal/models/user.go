package models

type User struct {
	Login    string `json:"login"`
	PassHash []byte `json:"passHash"`
	ID       int64  `json:"id"`
}

func NewUser(id int64, login string, passHash []byte) *User {
	return &User{
		ID:       id,
		Login:    login,
		PassHash: passHash,
	}
}
