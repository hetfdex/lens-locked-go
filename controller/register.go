package controller

import (
	"lens-locked-go/model"
	"lens-locked-go/service"
	"net/http"
)

type registerController struct {
	*controller
}

func NewRegisterController(us service.IUserService) *registerController {
	return &registerController{
		newController("/register", "view/register.gohtml", us),
	}
}

func (c *registerController) Post(w http.ResponseWriter, req *http.Request) {
	register := &model.RegisterForm{}

	err := parseForm(req, register)

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	err = register.Validate()

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	_, token, err := c.userService.Register(register)

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
