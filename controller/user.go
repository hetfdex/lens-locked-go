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

//Register user
func (c *userController) GetRegisterUser(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	parseSuccessFromRoute(req, data)

	c.registerView.Render(w, req, data)
}

func (c *userController) PostRegisterUser(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}
	form := &model.RegisterUser{}

	err := parseForm(req, form)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.registerView.Render(w, req, data)

		return
	}
	err = form.Validate()

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.registerView.Render(w, req, data)

		return
	}
	_, token, err := c.userService.Register(form)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.registerView.Render(w, req, data)

		return
	}
	cookie, err := makeCookie(token)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.registerView.Render(w, req, data)

		return
	}
	http.SetCookie(w, cookie)

	route := addSuccessToRoute(homeRoute, registerUserValue)

	util.Redirect(w, req, route)
}

//Login user
func (c *userController) GetLoginUser(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	parseSuccessFromRoute(req, data)

	c.loginView.Render(w, req, data)
}

func (c *userController) PostLoginUser(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}
	form := &model.LoginUser{}

	err := parseForm(req, form)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.loginView.Render(w, req, data)

		return
	}
	err = form.Validate()

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.loginView.Render(w, req, data)

		return
	}
	_, token, err := c.userService.LoginWithPassword(form)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.loginView.Render(w, req, data)

		return
	}
	cookie, err := makeCookie(token)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.loginView.Render(w, req, data)

		return
	}
	http.SetCookie(w, cookie)

	route := addSuccessToRoute(homeRoute, loginUserValue)

	util.Redirect(w, req, route)
}

//Logout user
func (c *userController) GetLogoutUser(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	user, err := context.User(req.Context())

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.logoutView.Render(w, req, data)

		return
	}
	err = c.userService.Logout(user)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.logoutView.Render(w, req, data)

		return
	}
	cookie, err := makeCookie(invalidTokenValue)

	if err != nil {
		data.Alert = err.Alert()

		w.WriteHeader(err.StatusCode)

		c.logoutView.Render(w, req, data)

		return
	}
	http.SetCookie(w, cookie)

	route := addSuccessToRoute(homeRoute, logoutUserValue)

	util.Redirect(w, req, route)
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
