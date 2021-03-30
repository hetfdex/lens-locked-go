package model

import "github.com/gofrs/uuid"

type Gallery struct {
	Base
	Name   string    `gorm:"not null"`
	UserId uuid.UUID `gorm:"not null;index"`
}
