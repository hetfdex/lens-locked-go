package model

import "lens-locked-go/util"

type UpdateForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *UpdateForm) Validate() *ApiError {
	if u.Name == "" {
		return NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("name"))
	}

	if u.Email == "" {
		return NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("email"))
	}

	if u.Password == "" {
		return NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}
