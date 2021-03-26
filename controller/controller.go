package controller

import (
	"errors"
	"lens-locked-go/model"
	"lens-locked-go/service"
	view2 "lens-locked-go/view"
	"log"
	"net/http"
)

type controller struct {
	Route       string
	view        *view2.View
	userService service.IUserService
}

func (c *controller) Get(w http.ResponseWriter, _ *http.Request) {
	data := &model.Data{}

	err := c.view.Render(w, data)

	if err != nil {
		c.handleError(w, err, data)
	}
}

func (c *controller) handleError(w http.ResponseWriter, err *model.ApiError, data *model.Data) {
	log.Println(err)

	data.Alert = err.Alert()

	c.view.Render(w, data)
}

func newController(route string, filename string, us service.IUserService) *controller {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	return &controller{
		Route:       route,
		view:        view2.New(filename),
		userService: us,
	}
}
