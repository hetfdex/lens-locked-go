package model

import "lens-locked-go/util"

type LoginUser struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *LoginUser) Validate() *Error {
	if u.Email == "" {
		return NewBadRequestApiError(util.MustNotBeEmptyErrorMessage("email"))
	}

	if u.Password == "" {
		return NewBadRequestApiError(util.MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}
