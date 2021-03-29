package model

import "github.com/gofrs/uuid"

type CreateGallery struct {
	Name string `schema:"name"`
}

func (c *CreateGallery) Validate() *Error {
	if c.Name == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("title"))
	}
	return nil
}

func (c *CreateGallery) Gallery(userId uuid.UUID) *Gallery {
	return &Gallery{
		Title:  c.Name,
		UserId: userId,
	}
}
