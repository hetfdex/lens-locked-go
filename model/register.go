package model

type Register struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
