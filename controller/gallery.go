package controller

import (
	"github.com/gorilla/mux"
	"lens-locked-go/context"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

const createGalleryRoute = "/gallery"
const createGalleryFilename = "view/gallery_create.gohtml"

const editGalleryRoute = "/gallery/{id}/edit"
const editGalleryFilename = "view/gallery_edit.gohtml"

const GalleryRouteName = "gallery"
const galleryRoute = "/gallery/{id}"
const galleryFilename = "view/gallery.gohtml"

const idKey = "id"

const invalidUUIDErrorMessage = "invalid gallery id"
const noUserInContextErrorMessage = "no user in context"

type galleryController struct {
	galleryView    *view.View
	createView     *view.View
	editView       *view.View
	router         *mux.Router
	galleryService service.IGalleryService
}

func NewGalleryController(r *mux.Router, gs service.IGalleryService) *galleryController {
	return &galleryController{
		galleryView:    view.New(galleryRoute, galleryFilename),
		createView:     view.New(createGalleryRoute, createGalleryFilename),
		editView:       view.New(editGalleryRoute, editGalleryFilename),
		router:         r,
		galleryService: gs,
	}
}

func (c *galleryController) GalleryGet(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}

	gallery, err := getGallery(req, c.galleryService)

	if err != nil {
		handleError(w, c.galleryView, err, data)

		return
	}
	data.Data = gallery

	c.galleryView.Render(w, data)
}

func (c *galleryController) CreateGet(w http.ResponseWriter, _ *http.Request) {
	data := &model.DataView{}

	c.createView.Render(w, data)
}

func (c *galleryController) CreatePost(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}
	create := &model.CreateGallery{}

	err := parseForm(req, create)

	if err != nil {
		handleError(w, c.createView, err, data)

		return
	}
	err = create.Validate()

	if err != nil {
		handleError(w, c.createView, err, data)

		return
	}
	user := context.User(req.Context())

	if user == nil {
		err = model.NewInternalServerApiError(noUserInContextErrorMessage)

		handleError(w, c.createView, err, data)

		return
	}
	gallery, err := c.galleryService.Create(create, user.ID)

	if err != nil {
		handleError(w, c.createView, err, data)

		return
	}
	url, err := makeUrl(c.router, GalleryRouteName, idKey, gallery.ID.String())

	if err != nil {
		handleError(w, c.createView, err, data)

		return
	}
	Redirect(w, req, url)
}

func (c *galleryController) EditGet(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}

	gallery, err := getGallery(req, c.galleryService)

	if err != nil {
		handleError(w, c.editView, err, data)

		return
	}
	data.Data = gallery

	c.editView.Render(w, data)
}

func (c *galleryController) EditPost(w http.ResponseWriter, req *http.Request) {
}

func (c *galleryController) GalleryRoute() string {
	return c.galleryView.Route()
}

func (c *galleryController) CreateRoute() string {
	return c.createView.Route()
}

func (c *galleryController) EditRoute() string {
	return c.editView.Route()
}
