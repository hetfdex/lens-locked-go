package controller

import (
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"lens-locked-go/context"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/util"
	"lens-locked-go/view"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

const invalidIdErrorMessage = "invalid gallery id"
const userNotOwnerErrorMessage = "galleries can only be edited by their owners"

type galleryController struct {
	indexGalleryView  *view.View
	createGalleryView *view.View
	editGalleryView   *view.View
	uploadGalleryView *view.View
	deleteGalleryView *view.View
	galleryView       *view.View
	galleryService    service.IGalleryService
	imageService      service.IImageService
}

func NewGalleryController(gs service.IGalleryService, is service.IImageService) *galleryController {
	return &galleryController{
		indexGalleryView:  view.New(indexGalleryRoute, indexGalleryFilename),
		createGalleryView: view.New(createGalleryRoute, createGalleryFilename),
		editGalleryView:   view.New(editGalleryRoute, editGalleryFilename),
		uploadGalleryView: view.New(uploadGalleryRoute, uploadGalleryFilename),
		deleteGalleryView: view.New(deleteGalleryRoute, galleryFilename),
		galleryView:       view.New(galleryRoute, galleryFilename),
		galleryService:    gs,
		imageService:      is,
	}
}

func (c *galleryController) IndexGet(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	user, err := context.User(req.Context())

	if err != nil {
		handleError(w, req, data, err, c.indexGalleryView)

		return
	}
	galleries, err := c.galleryService.GetAllByUserId(user.ID)

	if err != nil {
		handleError(w, req, data, err, c.indexGalleryView)

		return
	}
	data.Value = galleries

	c.indexGalleryView.Render(w, req, data)
}

func (c *galleryController) CreateGet(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	parseSuccessFromRoute(req, data)

	c.createGalleryView.Render(w, req, data)
}

func (c *galleryController) CreatePost(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}
	form := &model.CreateGallery{}

	err := parseForm(req, form)

	if err != nil {
		handleError(w, req, data, err, c.createGalleryView)

		return
	}
	err = form.Validate()

	if err != nil {
		handleError(w, req, data, err, c.createGalleryView)

		return
	}
	user, err := context.User(req.Context())

	if err != nil {
		handleError(w, req, data, err, c.createGalleryView)

		return
	}
	_, err = c.galleryService.Create(form, user.ID)

	if err != nil {
		handleError(w, req, data, err, c.createGalleryView)

		return
	}
	route := addSuccessToRoute(indexGalleryRoute, createGalleryValue)

	util.Redirect(w, req, route)
}

func (c *galleryController) EditGet(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	gallery, err := getGalleryWithPermission(req, c.galleryService)

	if err != nil {
		handleError(w, req, data, err, c.editGalleryView)

		return
	}
	data.Value = gallery

	c.editGalleryView.Render(w, req, data)
}

func (c *galleryController) EditPost(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}
	form := &model.EditGallery{}

	gallery, err := getGalleryWithPermission(req, c.galleryService)

	if err != nil {
		handleError(w, req, data, err, c.editGalleryView)

		return
	}
	err = parseForm(req, form)

	if err != nil {
		handleError(w, req, data, err, c.editGalleryView)

		return
	}
	err = form.Validate()

	if err != nil {
		handleError(w, req, data, err, c.editGalleryView)

		return
	}
	err = c.galleryService.Update(gallery, form)

	if err != nil {
		handleError(w, req, data, err, c.editGalleryView)

		return
	}
	route := addSuccessToRoute(indexGalleryRoute, editGalleryValue)

	util.Redirect(w, req, route)
}

func (c *galleryController) UploadGet(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	gallery, err := getGalleryWithPermission(req, c.galleryService)

	if err != nil {
		handleError(w, req, data, err, c.uploadGalleryView)

		return
	}
	data.Value = gallery

	c.uploadGalleryView.Render(w, req, data)
}

func (c *galleryController) UploadPost(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	gallery, err := getGalleryWithPermission(req, c.galleryService)

	if err != nil {
		handleError(w, req, data, err, c.uploadGalleryView)

		return
	}
	er := req.ParseMultipartForm(1 << 20) //1 megabyte

	if er != nil {
		handleError(w, req, data, model.NewInternalServerApiError(er.Error()), c.uploadGalleryView)

		return
	}
	fhs := req.MultipartForm.File["images"]

	for _, fh := range fhs {
		file, e := fh.Open()

		if e != nil {
			handleError(w, req, data, model.NewInternalServerApiError(e.Error()), c.uploadGalleryView)

			return
		}
		err = c.imageService.Create(file, fh.Filename, gallery.ID)

		if err != nil {
			handleError(w, req, data, err, c.uploadGalleryView)

			return
		}
		_ = file.Close()
	}
	route := addSuccessToRoute(indexGalleryRoute+"/"+gallery.ID.String(), uploadGalleryValue)

	util.Redirect(w, req, route)
}

func (c *galleryController) UploadDropboxPost(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	gallery, err := getGalleryWithPermission(req, c.galleryService)

	if err != nil {
		handleError(w, req, data, err, c.uploadGalleryView)

		return
	}
	er := req.ParseForm()

	if er != nil {
		handleError(w, req, data, model.NewInternalServerApiError(er.Error()), c.uploadGalleryView)

		return
	}
	urls := req.PostForm["urls"]

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)

		go downloadAndCreateDropboxImages(url, gallery.ID, &wg, c.imageService)
	}
	wg.Wait()
	route := addSuccessToRoute(indexGalleryRoute+"/"+gallery.ID.String(), uploadGalleryValue)

	util.Redirect(w, req, route)
}

func (c *galleryController) DeleteGet(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	gallery, err := getGalleryWithPermission(req, c.galleryService)

	if err != nil {
		handleError(w, req, data, err, c.deleteGalleryView)

		return
	}
	images, err := c.imageService.GetAllByGalleryId(gallery.ID)

	if err != nil {
		handleError(w, req, data, err, c.deleteGalleryView)

		return
	}

	for _, image := range images {
		err = c.imageService.Delete(image)

		if err != nil {
			handleError(w, req, data, err, c.deleteGalleryView)

			return
		}
	}
	err = c.galleryService.Delete(gallery)

	if err != nil {
		handleError(w, req, data, err, c.deleteGalleryView)

		return
	}
	route := addSuccessToRoute(indexGalleryRoute, deleteGalleryValue)

	util.Redirect(w, req, route)
}

func (c *galleryController) GalleryGet(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	parseSuccessFromRoute(req, data)

	gallery, err := getGallery(req, c.galleryService)

	if err != nil {
		handleError(w, req, data, err, c.galleryView)

		return
	}
	images, err := c.imageService.GetAllByGalleryId(gallery.ID)

	if err != nil {
		handleError(w, req, data, err, c.galleryView)

		return
	}
	gallery.Images = images

	data.Value = gallery

	c.galleryView.Render(w, req, data)
}

func (c *galleryController) IndexRoute() string {
	return c.indexGalleryView.Route()
}

func (c *galleryController) CreateRoute() string {
	return c.createGalleryView.Route()
}

func (c *galleryController) EditRoute() string {
	return c.editGalleryView.Route()
}

func (c *galleryController) UploadRoute() string {
	return c.uploadGalleryView.Route()
}

func (c *galleryController) UploadDropboxRoute() string {
	return uploadDropboxGalleryRoute
}

func (c *galleryController) DeleteRoute() string {
	return c.deleteGalleryView.Route()
}

func (c *galleryController) GalleryRoute() string {
	return c.galleryView.Route()
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

func downloadAndCreateDropboxImages(url string, galleryId uuid.UUID, wg *sync.WaitGroup, imageService service.IImageService) {
	defer wg.Done()

	res, e := http.Get(url)

	if e != nil {
		log.Println(url, e.Error())

		return
	}
	defer res.Body.Close()

	_, filename := filepath.Split(url)

	err := imageService.Create(res.Body, filename, galleryId)

	if err != nil {
		log.Println(url, err.Message)

		return
	}
}
