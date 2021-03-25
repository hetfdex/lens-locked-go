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
	if !register.Valid() {
		err = model.NewBadRequestApiError("invalid form")

		http.Error(w, err.Message, err.StatusCode)

		return
	}
	user := register.User()

	err = c.userService.Register(user)

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	cookie := makeCookie(user.Token)

	http.SetCookie(w, cookie)

	redirect(w, req, "/")
}
