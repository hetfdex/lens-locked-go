package model

type UserLogin struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (l *UserLogin) Validate() *Error {
	if l.Email == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("email"))
	}

	if l.Password == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}
