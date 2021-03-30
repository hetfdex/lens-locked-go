package controller

import (
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

const registerUserRoute = "/register"
const registerUserFilename = "view/user_register.gohtml"

const LoginUserRoute = "/login"
const loginUserFilename = "view/user_login.gohtml"

type userController struct {
	registerView *view.View
	loginView    *view.View
	userService  service.IUserService
}

func NewUserController(us service.IUserService) *userController {
	return &userController{
		registerView: view.New(registerUserRoute, registerUserFilename),
		loginView:    view.New(LoginUserRoute, loginUserFilename),
		userService:  us,
	}
}

//Register user
func (c *userController) GetRegisterUser(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}

	parseSuccessRoute(req, data)

	c.registerView.Render(w, data)
}

func (c *userController) PostRegisterUser(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}
	register := &model.RegisterUser{}

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

	route := makeSuccessRoute(homeRoute, registerUserValue)

	Redirect(w, req, route)
}

//Login user
func (c *userController) GetLoginUser(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}

	parseSuccessRoute(req, data)

	c.loginView.Render(w, data)
}

func (c *userController) PostLoginUser(w http.ResponseWriter, req *http.Request) {
	data := &model.DataView{}
	login := &model.LoginUser{}

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

	route := makeSuccessRoute(homeRoute, loginUserValue)

	Redirect(w, req, route)
}

func (c *userController) LoginUserRoute() string {
	return c.loginView.Route()
}

func (c *userController) RegisterUserRoute() string {
	return c.registerView.Route()
}
