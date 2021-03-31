package model

import "github.com/gofrs/uuid"

type Image struct {
	Base
	Source    string    `gorm:"not null"`
	Name      string    `gorm:"not null"`
	Extension string    `gorm:"not null"`
	GalleryId uuid.UUID `gorm:"not null;index"`
}
