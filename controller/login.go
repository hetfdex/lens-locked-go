package controller

import (
	"lens-locked-go/model"
	"lens-locked-go/service"
	"net/http"
)

type loginController struct {
	*controller
}

func NewLoginController(us *service.UserService) *loginController {
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
	if !login.Valid() {
		err = model.NewBadRequestApiError("invalid form")

		http.Error(w, err.Message, err.StatusCode)

		return
	}
	user, err := c.userService.LoginWithPassword(login)

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	err = c.userService.UpdateToken(user)

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	cookie := makeCookie(user.Token)

	http.SetCookie(w, cookie)

	redirect(w, req, "/")
}
