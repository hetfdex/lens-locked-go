package controller

import (
	"errors"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

type registerUserController struct {
	Route       string
	view        *view.View
	userService service.IUserService
}

func NewRegisterUserController(us service.IUserService) *registerUserController {
	return newRegisterUserController("/register", "view/user_register.gohtml", us)
}

func (c *registerUserController) Get(w http.ResponseWriter, _ *http.Request) {
	data := &model.DataView{}

	c.view.Render(w, data)
}

func (c *registerUserController) Post(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}
	register := &model.UserRegister{}

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

func newRegisterUserController(route string, filename string, us service.IUserService) *registerUserController {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	return &registerUserController{
		Route:       route,
		view:        view.New(filename),
		userService: us,
	}
}
