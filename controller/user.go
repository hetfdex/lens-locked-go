package controller

import (
	"lens-locked-go/context"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/util"
	"lens-locked-go/view"
	"net/http"
)

const registerUserRoute = "/register"
const registerUserFilename = "view/user_register.gohtml"

const loginUserRoute = "/login"
const loginUserFilename = "view/user_login.gohtml"

const logoutUserRoute = "/logout"

const invalidTokenValue = "invalidTokenValue"

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
		logoutView:   view.New(logoutUserRoute, homeFilename),
		userService:  us,
	}
}

func (c *userController) RegisterGet(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	parseSuccessFromRoute(req, data)

	c.registerView.Render(w, req, data)
}

func (c *userController) RegisterPost(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}
	form := &model.RegisterUser{}

	err := parseForm(req, form)

	if err != nil {
		handleError(w, req, data, err, c.registerView)

		return
	}
	err = form.Validate()

	if err != nil {
		handleError(w, req, data, err, c.registerView)

		return
	}
	_, token, err := c.userService.Register(form)

	if err != nil {
		handleError(w, req, data, err, c.registerView)

		return
	}
	cookie, err := makeCookie(token)

	if err != nil {
		handleError(w, req, data, err, c.registerView)

		return
	}
	http.SetCookie(w, cookie)

	route := addSuccessToRoute(homeRoute, registerUserValue)

	util.Redirect(w, req, route)
}

func (c *userController) LoginGet(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	parseSuccessFromRoute(req, data)

	c.loginView.Render(w, req, data)
}

func (c *userController) LoginPost(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}
	form := &model.LoginUser{}

	err := parseForm(req, form)

	if err != nil {
		handleError(w, req, data, err, c.loginView)

		return
	}
	err = form.Validate()

	if err != nil {
		handleError(w, req, data, err, c.loginView)

		return
	}
	_, token, err := c.userService.LoginWithPassword(form)

	if err != nil {
		handleError(w, req, data, err, c.loginView)

		return
	}
	cookie, err := makeCookie(token)

	if err != nil {
		handleError(w, req, data, err, c.loginView)

		return
	}
	http.SetCookie(w, cookie)

	route := addSuccessToRoute(homeRoute, loginUserValue)

	util.Redirect(w, req, route)
}

func (c *userController) LogoutGet(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	user, err := context.User(req.Context())

	if err != nil {
		handleError(w, req, data, err, c.logoutView)

		return
	}
	err = c.userService.Logout(user)

	if err != nil {
		handleError(w, req, data, err, c.logoutView)

		return
	}
	cookie, err := makeCookie(invalidTokenValue)

	if err != nil {
		handleError(w, req, data, err, c.logoutView)

		return
	}
	http.SetCookie(w, cookie)

	route := addSuccessToRoute(homeRoute, logoutUserValue)

	util.Redirect(w, req, route)
}

func (c *userController) RegisterRoute() string {
	return c.registerView.Route()
}

func (c *userController) LoginRoute() string {
	return c.loginView.Route()
}

func (c *userController) LogoutRoute() string {
	return c.logoutView.Route()
}

func makeCookie(value string) (*http.Cookie, *model.Error) {
	if value == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("value"))
	}
	return &http.Cookie{
		Name:     util.CookieName,
		Value:    value,
		HttpOnly: true,
	}, nil
}
