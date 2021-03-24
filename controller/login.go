package controller

import (
	"fmt"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"net/http"
)

type loginController struct {
	*controller
	*service.UserService
}

func NewLoginController(us *service.UserService) *loginController {
	return &loginController{
		newController("/login", "view/login.gohtml"),
		us,
	}
}

func (c *loginController) Login(w http.ResponseWriter, req *http.Request) {
	login := &model.LoginForm{}

	err := parseForm(req, login)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	user, err := c.Authenticate(login)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, _ = fmt.Fprintln(w, user)
}
