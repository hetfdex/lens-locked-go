package controller

import (
	"errors"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

type newGalleryController struct {
	Route          string
	view           *view.View
	galleryService service.IGalleryService
}

func NewCreateGalleryController(gs service.IGalleryService) *newGalleryController {
	return newCreateGalleryController("/gallery/create", "view/gallery_create.gohtml", gs)
}

func (c *newGalleryController) Get(w http.ResponseWriter, _ *http.Request) {
	data := &model.DataView{}

	c.view.Render(w, data)
}

func (c *newGalleryController) Post(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}
	create := &model.CreateGallery{}

	err := parseForm(req, create)

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	err = create.Validate()

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	_, err = c.galleryService.Create(create)

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	redirect(w, req, "/gallery")
}

func newCreateGalleryController(route string, filename string, gs service.IGalleryService) *newGalleryController {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	return &newGalleryController{
		Route:          route,
		view:           view.New(filename),
		galleryService: gs,
	}
}
