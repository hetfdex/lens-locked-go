package model

import "lens-locked-go/validator"

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (l *LoginForm) Validate() *ApiError {
	if validator.EmptyString(l.Email) {
		return NewInternalServerApiError("email must not be empty")
	}

	if validator.EmptyString(l.Password) {
		return NewInternalServerApiError("password must not be empty")
	}
	return nil
}
