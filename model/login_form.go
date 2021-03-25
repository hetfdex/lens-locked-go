package model

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (l *LoginForm) Validate() *ApiError {
	if l.Email == "" {
		return NewInternalServerApiError(MustNotBeEmptyErrorMessage("email"))
	}

	if l.Password == "" {
		return NewInternalServerApiError(MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}
