package controller

import (
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
	v, err := view.New(filename)

	if err != nil {
		panic(err)
	}
	return &controller{
		Route:       route,
		view:        v,
		userService: us,
	}
}
