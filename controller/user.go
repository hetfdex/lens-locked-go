package controller

import (
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

const LoginUserRoute = "/login"
const loginUserFilename = "view/user_login.gohtml"

const registerUserRoute = "/register"
const registerUserFilename = "view/user_register.gohtml"

type userController struct {
	loginView    *view.View
	registerView *view.View
	userService  service.IUserService
}

func NewUserController(us service.IUserService) *userController {
	return &userController{
		loginView:    view.New(LoginUserRoute, loginUserFilename),
		registerView: view.New(registerUserRoute, registerUserFilename),
		userService:  us,
	}
}

func (c *userController) RegisterGet(w http.ResponseWriter, _ *http.Request) {
	data := &model.DataView{}

	c.registerView.Render(w, data)
}

func (c *userController) LoginGet(w http.ResponseWriter, _ *http.Request) {
	data := &model.DataView{}

	c.loginView.Render(w, data)
}

func (c *userController) RegisterPost(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}
	register := &model.UserRegister{}

	err := parseForm(req, register)

	if err != nil {
		handleError(w, c.registerView, err, data)

		return
	}
	err = register.Validate()

	if err != nil {
		handleError(w, c.registerView, err, data)

		return
	}
	_, token, err := c.userService.Register(register)

	if err != nil {
		handleError(w, c.registerView, err, data)

		return
	}
	cookie, err := makeCookie(token)

	if err != nil {
		handleError(w, c.registerView, err, data)

		return
	}
	http.SetCookie(w, cookie)

	Redirect(w, req, homeRoute)
}

func (c *userController) LoginPost(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}
	login := &model.UserLogin{}

	err := parseForm(req, login)

	if err != nil {
		handleError(w, c.loginView, err, data)

		return
	}
	err = login.Validate()

	if err != nil {
		handleError(w, c.loginView, err, data)

		return
	}
	_, token, err := c.userService.LoginWithPassword(login)

	if err != nil {
		handleError(w, c.loginView, err, data)

		return
	}
	cookie, err := makeCookie(token)

	if err != nil {
		handleError(w, c.loginView, err, data)

		return
	}
	http.SetCookie(w, cookie)

	Redirect(w, req, homeRoute)
}

func (c *userController) LoginRoute() string {
	return c.loginView.Route()
}

func (c *userController) RegisterRoute() string {
	return c.registerView.Route()
}
