package controller

import (
	"errors"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"log"
	"net/http"
)

type controller struct {
	Route       string
	view        *model.View
	userService service.IUserService
}

func (c *controller) Get(w http.ResponseWriter, _ *http.Request) {
	data := &model.Data{}

	err := c.view.Render(w, data)

	if err != nil {
		c.handleError(w, err, data)
	}
}

func newController(route string, filename string, us service.IUserService) *controller {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	return &controller{
		Route:       route,
		view:        model.NewView(filename),
		userService: us,
	}
}

func (c *controller) handleError(w http.ResponseWriter, err *model.ApiError, data *model.Data) {
	log.Println(err)

	data.Alert = err.Alert()

	c.view.Render(w, data)
}
