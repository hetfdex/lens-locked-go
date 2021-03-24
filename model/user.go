package model

type User struct {
	Base
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Token        string `gorm:"-"`
	TokenHash    string `gorm:"not null;unique_index"`
}
