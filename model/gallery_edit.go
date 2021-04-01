package model

import "lens-locked-go/util"

type EditGallery struct {
	Name string `schema:"name"`
}

func (g *EditGallery) Validate() *Error {
	if g.Name == "" {
		return NewBadRequestApiError(util.MustNotBeEmptyErrorMessage("name"))
	}
	return nil
}
