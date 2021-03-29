package model

import "github.com/gofrs/uuid"

type EditGallery struct {
	Name string `schema:"name"`
}

func (g *EditGallery) Validate() *Error {
	if g.Name == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("title"))
	}
	return nil
}

func (g *EditGallery) Gallery(userId uuid.UUID) *Gallery {
	return &Gallery{
		Title:  g.Name,
		UserId: userId,
	}
}
