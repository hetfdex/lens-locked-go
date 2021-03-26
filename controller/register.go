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
	data := &model.DataView{}
	register := &model.RegisterView{}

	err := parseForm(req, register)

	if err != nil {
		c.handleError(w, err, data)

		return
	}
	err = register.Validate()

	if err != nil {
		c.handleError(w, err, data)

		return
	}
	_, token, err := c.userService.Register(register)

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
