package controller

import (
	"errors"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/util"
	"lens-locked-go/view"
	"net/http"
)

type controller struct {
	Route       string
	view        *view.View
	userService service.IUserService
}

func (c *controller) Get(w http.ResponseWriter, _ *http.Request) {
	alert := &model.Alert{
		Level:   "success",
		Message: "yay",
	}

	apiErr := c.view.Render(w, alert)

	if apiErr != nil {
		http.Error(w, apiErr.Message, apiErr.StatusCode)
	}
}

func newController(route string, filename string, us service.IUserService) *controller {
	if route == "" {
		panic(errors.New(util.MustNotBeEmptyErrorMessage("route")))
	}
	return &controller{
		Route:       route,
		view:        view.New(filename),
		userService: us,
	}
}
