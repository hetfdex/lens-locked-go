package model

import (
	"github.com/gofrs/uuid"
	"golang.org/x/oauth2"
)

type OAuth struct {
	Base
	UserId uuid.UUID `gorm:"not null;unique_index"`
	oauth2.Token
}