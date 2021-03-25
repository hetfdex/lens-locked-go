package model

type RegisterForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (r *RegisterForm) Validate() *ApiError {
	if r.Name == "" {
		return NewInternalServerApiError("name must not be empty")
	}

	if r.Email == "" {
		return NewInternalServerApiError("email must not be empty")
	}

	if r.Password == "" {
		return NewInternalServerApiError("password must not be empty")
	}
	return nil
}
