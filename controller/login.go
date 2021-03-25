package controller

import (
	"lens-locked-go/model"
	"lens-locked-go/service"
	"net/http"
)

type loginController struct {
	*controller
}

func NewLoginController(us service.IUserService) *loginController {
	return &loginController{
		newController("/login", "view/login.gohtml", us),
	}
}

func (c *loginController) Post(w http.ResponseWriter, req *http.Request) {
	login := &model.LoginForm{}

	err := parseForm(req, login)

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	err = login.Validate()

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	_, token, err := c.userService.LoginWithPassword(login)

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	cookie, err := makeCookie(token)

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	http.SetCookie(w, cookie)

	redirect(w, req, "/")
}
