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
	data := &model.DataView{}
	login := &model.LoginView{}

	err := parseForm(req, login)

	if err != nil {
		c.handleError(w, err, data)

		return
	}
	err = login.Validate()

	if err != nil {
		c.handleError(w, err, data)

		return
	}
	_, token, err := c.userService.LoginWithPassword(login)

	if err != nil {
		c.handleError(w, err, data)

		return
	}
	cookie, err := makeCookie(token)

	if err != nil {
		c.handleError(w, err, data)

		return
	}
	http.SetCookie(w, cookie)

	redirect(w, req, "/")
}
