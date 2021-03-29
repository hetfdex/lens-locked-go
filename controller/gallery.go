package controller

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

const idKey = "id"

const invalidUUIDErrorMessage = "invalid gallery id"

type galleryController struct {
	Route          string
	view           *view.View
	galleryService service.IGalleryService
}

func NewGalleryController(gs service.IGalleryService) *galleryController {
	return newGalleryController(galleryRoute, galleryFilename, gs)
}

func (c *galleryController) Get(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}

	vars := mux.Vars(req)

	id := uuid.FromStringOrNil(vars[idKey])

	if id == uuid.Nil {
		err := model.NewBadRequestApiError(invalidUUIDErrorMessage)

		handleError(c.view, w, err, data)

		return
	}
	gallery, err := c.galleryService.Get(id)

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	data.Data = gallery

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