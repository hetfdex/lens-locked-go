package model

type EditGallery struct {
	Name string `schema:"name"`
}

func (g *EditGallery) Validate() *Error {
	if g.Name == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("name"))
	}
	return nil
}
