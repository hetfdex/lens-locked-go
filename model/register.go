package model

type Register struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (r *Register) Validate() *ApiError {
	if r.Name == "" {
		return NewInternalServerApiError(MustNotBeEmptyErrorMessage("name"))
	}

	if r.Email == "" {
		return NewInternalServerApiError(MustNotBeEmptyErrorMessage("email"))
	}

	if r.Password == "" {
		return NewInternalServerApiError(MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}

func (r *Register) User(passwordHash string, tokenHash string) *User {
	return &User{
		Name:         r.Name,
		Email:        r.Email,
		PasswordHash: passwordHash,
		TokenHash:    tokenHash,
	}
}
