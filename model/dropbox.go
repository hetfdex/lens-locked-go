package model

import (
	"github.com/gofrs/uuid"
	"golang.org/x/oauth2"
)

type Dropbox struct {
	Base
	oauth2.Token
	UserId uuid.UUID `gorm:"not null;unique_index"`
}
