package model

import "lens-locked-go/util"

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (l *LoginForm) Validate() *ApiError {
	if l.Email == "" {
		return NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("email"))
	}

	if l.Password == "" {
		return NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}
