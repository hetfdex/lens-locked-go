package model

type RegisterForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (r *RegisterForm) Validate() *ApiError {
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
