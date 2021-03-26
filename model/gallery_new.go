package model

type NewGallery struct {
	Title string `schema:"title"`
}

func (g *NewGallery) Validate() *Error {
	if g.Title == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("title"))
	}
	return nil
}

func (g *NewGallery) Gallery() *Gallery {
	return &Gallery{
		Title: g.Title,
	}
}
