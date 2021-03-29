package model

import "github.com/gofrs/uuid"

type CreateGallery struct {
	Name string `schema:"name"`
}

func (g *CreateGallery) Validate() *Error {
	if g.Name == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("title"))
	}
	return nil
}

func (g *CreateGallery) Gallery(userId uuid.UUID) *Gallery {
	return &Gallery{
		Title:  g.Name,
		UserId: userId,
	}
}
