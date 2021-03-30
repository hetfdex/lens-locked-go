package controller

import (
	"lens-locked-go/context"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

const registerUserRoute = "/register"
const registerUserFilename = "view/user_register.gohtml"

const loginUserRoute = "/login"
const loginUserFilename = "view/user_login.gohtml"

const logoutUserRoute = "/logout"
const logoutUserFilename = homeFilename

type userController struct {
	registerView *view.View
	loginView    *view.View
	logoutView   *view.View
	userService  service.IUserService
}

func NewUserController(us service.IUserService) *userController {
	return &userController{
		registerView: view.New(registerUserRoute, registerUserFilename),
		loginView:    view.New(loginUserRoute, loginUserFilename),
		logoutView:   view.New(logoutUserRoute, logoutUserFilename),
		userService:  us,
	}
}

//Register user
func (c *userController) GetRegisterUser(w http.ResponseWriter, req *http.Request) {
	viewData := &model.DataView{}

	parseSuccessRoute(req, viewData)

	c.registerView.Render(w, viewData)
}

func (c *userController) PostRegisterUser(w http.ResponseWriter, req *http.Request) {
	viewData := &model.DataView{}
	register := &model.RegisterUser{}

	err := parseForm(req, register)

	if err != nil {
		handleError(w, c.registerView, err, viewData)

		return
	}
	err = register.Validate()

	if err != nil {
		handleError(w, c.registerView, err, viewData)

		return
	}
	_, token, err := c.userService.Register(register)

	if err != nil {
		handleError(w, c.registerView, err, viewData)

		return
	}
	cookie, err := makeCookie(token)

	if err != nil {
		handleError(w, c.registerView, err, viewData)

		return
	}
	http.SetCookie(w, cookie)

	route := makeSuccessRoute(homeRoute, registerUserValue)

	Redirect(w, req, route)
}

//Login user
func (c *userController) GetLoginUser(w http.ResponseWriter, req *http.Request) {
	viewData := &model.DataView{}

	parseSuccessRoute(req, viewData)

	c.loginView.Render(w, viewData)
}

func (c *userController) PostLoginUser(w http.ResponseWriter, req *http.Request) {
	viewData := &model.DataView{}
	login := &model.LoginUser{}

	err := parseForm(req, login)

	if err != nil {
		handleError(w, c.loginView, err, viewData)

		return
	}
	err = login.Validate()

	if err != nil {
		handleError(w, c.loginView, err, viewData)

		return
	}
	_, token, err := c.userService.LoginWithPassword(login)

	if err != nil {
		handleError(w, c.loginView, err, viewData)

		return
	}
	cookie, err := makeCookie(token)

	if err != nil {
		handleError(w, c.loginView, err, viewData)

		return
	}
	http.SetCookie(w, cookie)

	route := makeSuccessRoute(homeRoute, loginUserValue)

	Redirect(w, req, route)
}

//Logout user
func (c *userController) GetLogoutUser(w http.ResponseWriter, req *http.Request) {
	viewData := &model.DataView{}

	user, err := context.User(req.Context())

	if err != nil {
		handleError(w, c.logoutView, err, viewData)

		return
	}
	err = c.userService.Logout(user)

	if err != nil {
		handleError(w, c.logoutView, err, viewData)

		return
	}
	cookie, err := makeCookie(invalidTokenValue)

	if err != nil {
		handleError(w, c.logoutView, err, viewData)

		return
	}
	http.SetCookie(w, cookie)

	route := makeSuccessRoute(homeRoute, logoutUserValue)

	Redirect(w, req, route)
}

func RegisterUserRoute() string {
	return registerUserRoute
}

func LoginUserRoute() string {
	return loginUserRoute
}

func LogoutUserRoute() string {
	return logoutUserRoute
}
