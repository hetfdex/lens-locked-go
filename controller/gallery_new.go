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

func NewNewGalleryController(gs service.IGalleryService) *newGalleryController {
	return newNewGalleryController("/gallery/new", "view/gallery_new.gohtml", gs)
}

func (c *newGalleryController) Get(w http.ResponseWriter, _ *http.Request) {
	data := &model.DataView{}

	c.view.Render(w, data)
}

func (c *newGalleryController) Post(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}
	newG := &model.NewGallery{}

	err := parseForm(req, newG)

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	err = newG.Validate()

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	_, err = c.galleryService.New(newG)

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	redirect(w, req, "/gallery")
}

func newNewGalleryController(route string, filename string, gs service.IGalleryService) *newGalleryController {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	return &newGalleryController{
		Route:          route,
		view:           view.New(filename),
		galleryService: gs,
	}
}
