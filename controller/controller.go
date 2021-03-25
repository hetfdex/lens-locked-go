package controller

import (
	"errors"
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
	apiErr := c.view.Render(w, nil)

	if apiErr != nil {
		http.Error(w, apiErr.Message, apiErr.StatusCode)
	}
}

func newController(route string, filename string, us service.IUserService) *controller {
	if route == "" {
		panic(errors.New("route must not be empty"))
	}
	return &controller{
		Route:       route,
		view:        view.New(filename),
		userService: us,
	}
}
