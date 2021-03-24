package controller

import (
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

type controller struct {
	Route       string
	view        *view.View
	userService *service.UserService
}

func (c *controller) Get(w http.ResponseWriter, _ *http.Request) {
	err := c.view.Render(w, nil)

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)
	}
}

func newController(route string, filename string, us *service.UserService) *controller {
	return &controller{
		Route:       route,
		view:        view.New(filename),
		userService: us,
	}
}
