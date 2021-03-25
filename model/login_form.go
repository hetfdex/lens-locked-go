package model

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (l *LoginForm) Validate() *ApiError {
	if l.Email == "" {
		return NewInternalServerApiError("email must not be empty")
	}

	if l.Password == "" {
		return NewInternalServerApiError("password must not be empty")
	}
	return nil
}
