package models

type User struct {
	Login    string
	PassHash []byte
	Email    string
}
