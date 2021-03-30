package model

type RegisterUser struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *RegisterUser) Validate() *Error {
	if u.Name == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("name"))
	}

	if u.Email == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("email"))
	}

	if u.Password == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}
