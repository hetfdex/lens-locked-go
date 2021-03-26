package model

type Update struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *Update) Validate() *ApiError {
	if u.Name == "" {
		return NewInternalServerApiError(MustNotBeEmptyErrorMessage("name"))
	}

	if u.Email == "" {
		return NewInternalServerApiError(MustNotBeEmptyErrorMessage("email"))
	}

	if u.Password == "" {
		return NewInternalServerApiError(MustNotBeEmptyErrorMessage("password"))
	}
	return nil
}
