package controller

import (
	"errors"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

type createGalleryController struct {
	Route          string
	view           *view.View
	galleryService service.IGalleryService
}

func NewCreateGalleryController(gs service.IGalleryService) *createGalleryController {
	return newCreateGalleryController("/gallery/create", "view/gallery_create.gohtml", gs)
}

func (c *createGalleryController) Get(w http.ResponseWriter, _ *http.Request) {
	data := &model.DataView{}

	c.view.Render(w, data)
}

func newCreateGalleryController(route string, filename string, gs service.IGalleryService) *createGalleryController {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	return &createGalleryController{
		Route:          route,
		view:           view.New(filename),
		galleryService: gs,
	}
}
