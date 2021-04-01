package model

import "lens-locked-go/util"

type CreateGallery struct {
	Name string `schema:"name"`
}

func (g *CreateGallery) Validate() *Error {
	if g.Name == "" {
		return NewBadRequestApiError(util.MustNotBeEmptyErrorMessage("name"))
	}
	return nil
}
