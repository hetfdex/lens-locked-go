package model

type LoginView struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (l *LoginView) Validate() *Error {
	if l.Email == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("email"))
	}

	if l.Password == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}
