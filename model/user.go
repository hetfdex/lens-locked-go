package model

type User struct {
	Base
	Name         string
	Email        string `gorm:"not null;unique_index"`
	PasswordHash string `gorm:"not null"`
	Token        string `gorm:"-"`
	TokenHash    string `gorm:"not null;unique_index"`
}

func NewUserFromRegister(register *RegisterForm, passwordHash string, token string, tokenHash string) *User {
	return &User{
		Name:         register.Name,
		Email:        register.Email,
		PasswordHash: passwordHash,
		Token:        token,
		TokenHash:    tokenHash,
	}
}
