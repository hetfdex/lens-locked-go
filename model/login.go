package model

import "lens-locked-go/util"

type Login struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (l *Login) Validate() *ApiError {
	if l.Email == "" {
		return NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("email"))
	}

	if l.Password == "" {
		return NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}
