package model

type CreateGallery struct {
	Title string `schema:"title"`
}

func (c *CreateGallery) Validate() *Error {
	if c.Title == "" {
		return NewBadRequestApiError(MustNotBeEmptyErrorMessage("title"))
	}
	return nil
}

func (c *CreateGallery) Gallery() *Gallery {
	return &Gallery{
		Title: c.Title,
	}
}
