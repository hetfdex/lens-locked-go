package controller

import (
	"errors"
	"lens-locked-go/context"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

const noUserInContextErrorMessage = "no user found in context"

type createGalleryController struct {
	Route          string
	view           *view.View
	galleryService service.IGalleryService
}

func NewCreateGalleryController(gs service.IGalleryService) *createGalleryController {
	return newCreateGalleryController(createGalleryRoute, createGalleryFilename, gs)
}

func (c *createGalleryController) Get(w http.ResponseWriter, _ *http.Request) {
	data := &model.DataView{}

	c.view.Render(w, data)
}

func (c *createGalleryController) Post(w http.ResponseWriter, req *http.Request) {
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
	user := context.User(req.Context())

	if user == nil {
		err = model.NewInternalServerApiError(noUserInContextErrorMessage)

		handleError(c.view, w, err, data)

		return
	}
	gallery, err := c.galleryService.Create(create, user.ID)

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	redirectToGallery(w, req, gallery.ID)
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
