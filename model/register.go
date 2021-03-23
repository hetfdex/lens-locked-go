package model

type Register struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (r *Register) User() *User {
	return &User{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}
