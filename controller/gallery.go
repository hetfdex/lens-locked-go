package controller

import (
	"errors"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

type galleryController struct {
	Route          string
	view           *view.View
	galleryService service.IGalleryService
}

func NewGalleryController(gs service.IGalleryService) *galleryController {
	return newGalleryController("/gallery", "view/gallery.gohtml", gs)
}

func (c *galleryController) Get(w http.ResponseWriter, _ *http.Request) {
	data := &model.DataView{}

	c.view.Render(w, data)
}

func newGalleryController(route string, filename string, gs service.IGalleryService) *galleryController {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	return &galleryController{
		Route:          route,
		view:           view.New(filename),
		galleryService: gs,
	}
}
