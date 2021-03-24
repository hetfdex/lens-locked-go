package controller

import (
	"lens-locked-go/model"
	"lens-locked-go/service"
	"net/http"
)

type registerController struct {
	*controller
	*service.UserService
}

func NewRegisterController(us *service.UserService) *registerController {
	return &registerController{
		newController("/register", "view/register.gohtml"),
		us,
	}
}

func (c *registerController) Register(w http.ResponseWriter, req *http.Request) {
	register := &model.RegisterForm{}

	err := parseForm(req, register)

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	if !validRegisterForm(register) {
		err = model.NewBadRequestApiError("invalid form")

		http.Error(w, err.Message, err.StatusCode)

		return
	}
	user := register.User()

	err = c.Create(user)

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	cookie := makeCookie(user.Email)

	http.SetCookie(w, cookie)

	redirect(w, req, "/")
}
