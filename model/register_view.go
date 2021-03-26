package model

type RegisterView struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (r *RegisterView) Validate() *Error {
	if r.Name == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("name"))
	}

	if r.Email == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("email"))
	}

	if r.Password == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}

func (r *RegisterView) User(passwordHash string, tokenHash string) *User {
	return &User{
		Name:         r.Name,
		Email:        r.Email,
		PasswordHash: passwordHash,
		TokenHash:    tokenHash,
	}
}
