package model

import "lens-locked-go/util"

type RegisterUser struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *RegisterUser) Validate() *Error {
	if u.Name == "" {
		return NewBadRequestApiError(util.MustNotBeEmptyErrorMessage("name"))
	}

	if u.Email == "" {
		return NewBadRequestApiError(util.MustNotBeEmptyErrorMessage("email"))
	}

	if u.Password == "" {
		return NewBadRequestApiError(util.MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}
