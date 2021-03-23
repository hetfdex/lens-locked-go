package model

type User struct {
	Base
	Name  string
	Email string `gorm:"not null;unique_index"`
}
