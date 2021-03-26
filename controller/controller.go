package controller

import (
	"errors"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/util"
	"net/http"
)

type controller struct {
	Route       string
	view        *model.View
	userService service.IUserService
}

func (c *controller) Get(w http.ResponseWriter, _ *http.Request) {
	apiErr := c.view.Render(w, nil)

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
		view:        model.NewView(filename),
		userService: us,
	}
}
