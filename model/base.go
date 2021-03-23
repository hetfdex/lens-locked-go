package model

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (b *Base) BeforeCreate(*gorm.DB) error {
	id, err := uuid.NewV4()

	if err != nil {
		return err
	}
	b.ID = id

	return nil
}
