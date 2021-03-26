package controller

import (
	"errors"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

type registerController struct {
	Route       string
	view        *view.View
	userService service.IUserService
}

func NewRegisterController(us service.IUserService) *registerController {
	return newRegisterController("/register", "view/register.gohtml", us)
}

func (c *registerController) Get(w http.ResponseWriter, _ *http.Request) {
	data := &model.DataView{}

	c.view.Render(w, data)
}

func (c *registerController) Post(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}
	register := &model.RegisterView{}

	err := parseForm(req, register)

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	err = register.Validate()

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	_, token, err := c.userService.Register(register)

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	cookie, err := makeCookie(token)

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	http.SetCookie(w, cookie)

	redirect(w, req, "/")
}

func newRegisterController(route string, filename string, us service.IUserService) *registerController {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	return &registerController{
		Route:       route,
		view:        view.New(filename),
		userService: us,
	}
}
