package model

type Login struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (l *Login) Validate() *ApiError {
	if l.Email == "" {
		return NewInternalServerApiError(MustNotBeEmptyErrorMessage("email"))
	}

	if l.Password == "" {
		return NewInternalServerApiError(MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}
