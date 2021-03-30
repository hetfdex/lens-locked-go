package controller

import (
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"lens-locked-go/context"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/util"
	"lens-locked-go/view"
	"net/http"
)

const indexGalleryRoute = "/gallery"
const indexGalleryFilename = "view/gallery_index.gohtml"

const createGalleryRoute = "/gallery/create"
const createGalleryFilename = "view/gallery_create.gohtml"

const editGalleryRoute = "/gallery/{gallery_id}/edit"
const editGalleryFilename = "view/gallery_edit.gohtml"

const deleteGalleryRoute = "/gallery/{gallery_id}/delete"

const galleryRoute = "/gallery/{gallery_id}"
const galleryFilename = "view/gallery.gohtml"

const galleryIdKey = "gallery_id"

const invalidIdErrorMessage = "invalid gallery id"
const userNotOwnerErrorMessage = "galleries can only be edited by their owners"

type galleryController struct {
	indexGalleryView  *view.View
	createGalleryView *view.View
	editGalleryView   *view.View
	deleteGalleryView *view.View
	galleryView       *view.View
	galleryService    service.IGalleryService
}

func NewGalleryController(gs service.IGalleryService) *galleryController {
	return &galleryController{
		indexGalleryView:  view.New(indexGalleryRoute, indexGalleryFilename),
		createGalleryView: view.New(createGalleryRoute, createGalleryFilename),
		editGalleryView:   view.New(editGalleryRoute, editGalleryFilename),
		deleteGalleryView: view.New(deleteGalleryRoute, galleryFilename),
		galleryView:       view.New(galleryRoute, galleryFilename),
		galleryService:    gs,
	}
}

//Index gallery
func (c *galleryController) GetIndexGallery(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	user, err := context.User(req.Context())

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.indexGalleryView.Render(w, req, data)

		return
	}
	galleries, err := c.galleryService.GetAllByUserId(user.ID)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.indexGalleryView.Render(w, req, data)

		return
	}
	data.User = user
	data.Value = galleries

	c.indexGalleryView.Render(w, req, data)
}

//Create gallery
func (c *galleryController) GetCreateGallery(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	parseSuccessFromRoute(req, data)

	c.createGalleryView.Render(w, req, data)
}

func (c *galleryController) PostCreateGallery(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}
	form := &model.CreateGallery{}

	err := parseForm(req, form)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.createGalleryView.Render(w, req, data)

		return
	}
	err = form.Validate()

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.createGalleryView.Render(w, req, data)

		return
	}
	user, err := context.User(req.Context())

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.createGalleryView.Render(w, req, data)

		return
	}
	gallery, err := c.galleryService.Create(form, user.ID)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.createGalleryView.Render(w, req, data)

		return
	}
	route := makeRouteFromId(gallery.ID)

	route = addSuccessToRoute(route, createGalleryValue)

	util.Redirect(w, req, route)
}

//Edit gallery
func (c *galleryController) GetEditGallery(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	gallery, err := getGalleryWithPermission(req, c.galleryService)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.editGalleryView.Render(w, req, data)

		return
	}
	data.Value = gallery

	c.editGalleryView.Render(w, req, data)
}

func (c *galleryController) PostEditGallery(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}
	form := &model.EditGallery{}

	gallery, err := getGalleryWithPermission(req, c.galleryService)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.editGalleryView.Render(w, req, data)

		return
	}
	err = parseForm(req, form)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.editGalleryView.Render(w, req, data)

		return
	}
	err = form.Validate()

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.editGalleryView.Render(w, req, data)

		return
	}
	err = c.galleryService.Edit(gallery, form)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.editGalleryView.Render(w, req, data)

		return
	}
	route := makeRouteFromId(gallery.ID)

	route = addSuccessToRoute(route, editGalleryValue)

	util.Redirect(w, req, route)
}

//Delete gallery
func (c *galleryController) GetDeleteGallery(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	gallery, err := getGalleryWithPermission(req, c.galleryService)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.deleteGalleryView.Render(w, req, data)

		return
	}
	err = c.galleryService.Delete(gallery)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.deleteGalleryView.Render(w, req, data)

		return
	}
	route := addSuccessToRoute(indexGalleryRoute, deleteGalleryValue)

	util.Redirect(w, req, route)
}

//Gallery
func (c *galleryController) GetGallery(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	parseSuccessFromRoute(req, data)

	gallery, err := getGallery(req, c.galleryService)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.galleryView.Render(w, req, data)

		return
	}
	data.Value = gallery

	c.galleryView.Render(w, req, data)
}

func IndexGalleryRoute() string {
	return indexGalleryRoute
}

func CreateGalleryRoute() string {
	return createGalleryRoute
}

func EditGalleryRoute() string {
	return editGalleryRoute
}

func DeleteGalleryRoute() string {
	return deleteGalleryRoute
}

func GalleryRoute() string {
	return galleryRoute
}

func makeRouteFromId(id uuid.UUID) string {
	return indexGalleryRoute + "/" + id.String()
}

func getGallery(req *http.Request, gs service.IGalleryService) (*model.Gallery, *model.Error) {
	vars := mux.Vars(req)

	id := uuid.FromStringOrNil(vars[galleryIdKey])

	if id == uuid.Nil {
		err := model.NewBadRequestApiError(invalidIdErrorMessage)

		return nil, err
	}
	gallery, err := gs.GetById(id)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func getGalleryWithPermission(req *http.Request, gs service.IGalleryService) (*model.Gallery, *model.Error) {
	gallery, err := getGallery(req, gs)

	if err != nil {
		return nil, err
	}
	user, err := context.User(req.Context())

	if err != nil {
		return nil, err
	}

	if gallery.UserId != user.ID {
		err = model.NewForbiddenApiError(userNotOwnerErrorMessage)

		return nil, err
	}
	return gallery, nil
}
