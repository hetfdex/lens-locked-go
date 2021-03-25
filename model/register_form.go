package model

import "lens-locked-go/util"

type RegisterForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (r *RegisterForm) Validate() *ApiError {
	if r.Name == "" {
		return NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("name"))
	}

	if r.Email == "" {
		return NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("email"))
	}

	if r.Password == "" {
		return NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}

func (r *RegisterForm) User(passwordHash string, tokenHash string) *User {
	return &User{
		Name:         r.Name,
		Email:        r.Email,
		PasswordHash: passwordHash,
		TokenHash:    tokenHash,
	}
}
