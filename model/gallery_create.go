package model

type CreateGallery struct {
	Name string `schema:"name"`
}

func (g *CreateGallery) Validate() *Error {
	if g.Name == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("name"))
	}
	return nil
}
