package controller

import (
	"errors"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

type controller struct {
	Route       string
	view        *view.View
	userService service.IUserService
}

func (c *controller) Get(w http.ResponseWriter, _ *http.Request) {
	data := &model.Data{}

	c.view.Render(w, data)
}

func (c *controller) handleError(w http.ResponseWriter, err *model.ApiError, data *model.Data) {
	data.Alert = err.Alert()

	w.WriteHeader(err.StatusCode)

	c.view.Render(w, data)
}

func newController(route string, filename string, us service.IUserService) *controller {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	return &controller{
		Route:       route,
		view:        view.New(filename),
		userService: us,
	}
}
