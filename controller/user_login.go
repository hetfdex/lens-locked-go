package controller

import (
	"errors"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

type loginController struct {
	Route       string
	view        *view.View
	userService service.IUserService
}

func NewLoginUserController(us service.IUserService) *loginController {
	return newLoginUserController("/login", "view/user_login.gohtml", us)
}

func (c *loginController) Get(w http.ResponseWriter, _ *http.Request) {
	data := &model.DataView{}

	c.view.Render(w, data)
}

func (c *loginController) Post(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}
	login := &model.LoginView{}

	err := parseForm(req, login)

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	err = login.Validate()

	if err != nil {
		handleError(c.view, w, err, data)

		return
	}
	_, token, err := c.userService.LoginWithPassword(login)

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

func newLoginUserController(route string, filename string, us service.IUserService) *loginController {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	return &loginController{
		Route:       route,
		view:        view.New(filename),
		userService: us,
	}
}
