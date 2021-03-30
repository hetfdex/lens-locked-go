package model

import "github.com/gofrs/uuid"

type Gallery struct {
	Base
	Title  string    `gorm:"not null"`
	UserId uuid.UUID `gorm:"not null;index"`
}
