package model

import "lens-locked-go/validator"

type RegisterForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (r *RegisterForm) Validate() *ApiError {
	if validator.EmptyString(r.Name) {
		return NewInternalServerApiError("name must not be empty")
	}

	if validator.EmptyString(r.Email) {
		return NewInternalServerApiError("email must not be empty")
	}

	if validator.EmptyString(r.Password) {
		return NewInternalServerApiError("password must not be empty")
	}
	return nil
}
