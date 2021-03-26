package model

import "github.com/gofrs/uuid"

type Gallery struct {
	Base
	UserId uuid.UUID `gorm:"not null;index"`
	Title  string    `gorm:"not null"`
}
