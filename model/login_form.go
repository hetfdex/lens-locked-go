package model

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (l *LoginForm) User() *User {
	return &User{
		Email:    l.Email,
		Password: l.Password,
	}
}

func (l *LoginForm) Valid() bool {
	if l.Email == "" || l.Password == "" {
		return false
	}
	return true
}
