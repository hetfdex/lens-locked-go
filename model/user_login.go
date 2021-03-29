package model

type LoginUser struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *LoginUser) Validate() *Error {
	if u.Email == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("email"))
	}

	if u.Password == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}
