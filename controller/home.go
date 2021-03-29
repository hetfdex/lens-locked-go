package controller

import (
	"errors"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

type homeController struct {
	Route       string
	view        *view.View
	userService service.IUserService
}

func NewHomeController(us service.IUserService) *homeController {
	return newHomeController(homeRoute, homeFilename, us)
}

func (c *homeController) Get(w http.ResponseWriter, _ *http.Request) {
	data := &model.DataView{}

	c.view.Render(w, data)
}

func (c *homeController) handleError(w http.ResponseWriter, err *model.Error, data *model.DataView) {
	data.Alert = err.Alert()

	w.WriteHeader(err.StatusCode)

	c.view.Render(w, data)
}

func newHomeController(route string, filename string, us service.IUserService) *homeController {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	return &homeController{
		Route:       route,
		view:        view.New(filename),
		userService: us,
	}
}
