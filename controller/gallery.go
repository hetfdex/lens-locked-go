package controller

import (
	"github.com/gorilla/mux"
	"lens-locked-go/context"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

const indexGalleryRoute = "/gallery"
const indexGalleryFilename = "view/gallery_index.gohtml"

const createGalleryRoute = "/gallery/create"
const createGalleryFilename = "view/gallery_create.gohtml"

const editGalleryRoute = "/gallery/{id}/edit"
const editGalleryFilename = "view/gallery_edit.gohtml"

const deleteGalleryRoute = "/gallery/{id}/delete"

const GalleryRouteName = "gallery"
const galleryRoute = "/gallery/{id}"
const galleryFilename = "view/gallery.gohtml"

const idKey = "id"

const invalidUUIDErrorMessage = "invalid gallery id"
const userNotOwnerErrorMessage = "galleries can only be edited by their owners"

type galleryController struct {
	indexGalleryView  *view.View
	createGalleryView *view.View
	editGalleryView   *view.View
	deleteGalleryView *view.View
	galleryView       *view.View
	router            *mux.Router
	galleryService    service.IGalleryService
}

func NewGalleryController(r *mux.Router, gs service.IGalleryService) *galleryController {
	return &galleryController{
		indexGalleryView:  view.New(indexGalleryRoute, indexGalleryFilename),
		createGalleryView: view.New(createGalleryRoute, createGalleryFilename),
		editGalleryView:   view.New(editGalleryRoute, editGalleryFilename),
		deleteGalleryView: view.New(deleteGalleryRoute, galleryFilename),
		galleryView:       view.New(galleryRoute, galleryFilename),
		router:            r,
		galleryService:    gs,
	}
}

//Index gallery
func (c *galleryController) GetIndexGallery(w http.ResponseWriter, req *http.Request) {
	viewData := &model.DataView{}

	user, err := context.User(req.Context())

	if err != nil {
		handleError(w, c.indexGalleryView, err, viewData)

		return
	}
	galleries, err := c.galleryService.GetAllByUserId(user.ID)

	if err != nil {
		handleError(w, c.indexGalleryView, err, viewData)

		return
	}
	viewData.Data = galleries

	c.indexGalleryView.Render(w, viewData)
}

//Create gallery
func (c *galleryController) GetCreateGallery(w http.ResponseWriter, req *http.Request) {
	viewData := &model.DataView{}

	parseSuccessRoute(req, viewData)

	c.createGalleryView.Render(w, viewData)
}

func (c *galleryController) PostCreateGallery(w http.ResponseWriter, req *http.Request) {
	viewData := &model.DataView{}
	create := &model.CreateGallery{}

	err := parseForm(req, create)

	if err != nil {
		handleError(w, c.createGalleryView, err, viewData)

		return
	}
	err = create.Validate()

	if err != nil {
		handleError(w, c.createGalleryView, err, viewData)

		return
	}
	user, err := context.User(req.Context())

	if err != nil {
		handleError(w, c.createGalleryView, err, viewData)

		return
	}
	gallery, err := c.galleryService.Create(create, user.ID)

	if err != nil {
		handleError(w, c.createGalleryView, err, viewData)

		return
	}
	url, err := makeUrl(c.router, GalleryRouteName, idKey, gallery.ID.String())

	if err != nil {
		handleError(w, c.createGalleryView, err, viewData)

		return
	}
	route := makeSuccessRoute(url, createGalleryValue)

	Redirect(w, req, route)
}

//Edit gallery
func (c *galleryController) GetEditGallery(w http.ResponseWriter, req *http.Request) {
	viewData := &model.DataView{}

	gallery, err := getGalleryWithPermission(req, c.galleryService)

	if err != nil {
		handleError(w, c.editGalleryView, err, viewData)

		return
	}
	viewData.Data = gallery

	c.editGalleryView.Render(w, viewData)
}

func (c *galleryController) PostEditGallery(w http.ResponseWriter, req *http.Request) {
	viewData := &model.DataView{}
	edit := &model.EditGallery{}

	gallery, err := getGalleryWithPermission(req, c.galleryService)

	if err != nil {
		handleError(w, c.editGalleryView, err, viewData)

		return
	}
	err = parseForm(req, edit)

	if err != nil {
		handleError(w, c.editGalleryView, err, viewData)

		return
	}
	err = edit.Validate()

	if err != nil {
		handleError(w, c.editGalleryView, err, viewData)

		return
	}
	gallery, err = c.galleryService.Edit(gallery, edit)

	if err != nil {
		handleError(w, c.editGalleryView, err, viewData)

		return
	}
	url, err := makeUrl(c.router, GalleryRouteName, idKey, gallery.ID.String())

	if err != nil {
		handleError(w, c.editGalleryView, err, viewData)

		return
	}
	route := makeSuccessRoute(url, editGalleryValue)

	Redirect(w, req, route)
}

//Delete gallery
func (c *galleryController) GetDeleteGallery(w http.ResponseWriter, req *http.Request) {
	viewData := &model.DataView{}

	gallery, err := getGalleryWithPermission(req, c.galleryService)

	if err != nil {
		handleError(w, c.deleteGalleryView, err, viewData)

		return
	}
	err = c.galleryService.Delete(gallery)

	if err != nil {
		handleError(w, c.deleteGalleryView, err, viewData)

		return
	}
	route := makeSuccessRoute(homeRoute, deleteGalleryValue)

	Redirect(w, req, route)
}

//Gallery
func (c *galleryController) GetGallery(w http.ResponseWriter, req *http.Request) {
	viewData := &model.DataView{}

	parseSuccessRoute(req, viewData)

	gallery, err := getGallery(req, c.galleryService)

	if err != nil {
		handleError(w, c.galleryView, err, viewData)

		return
	}
	viewData.Data = gallery

	c.galleryView.Render(w, viewData)
}

func (c *galleryController) IndexGalleryRoute() string {
	return c.indexGalleryView.Route()
}

func (c *galleryController) CreateGalleryRoute() string {
	return c.createGalleryView.Route()
}

func (c *galleryController) EditGalleryRoute() string {
	return c.editGalleryView.Route()
}

func (c *galleryController) DeleteGalleryRoute() string {
	return c.deleteGalleryView.Route()
}

func (c *galleryController) GalleryRoute() string {
	return c.galleryView.Route()
}
