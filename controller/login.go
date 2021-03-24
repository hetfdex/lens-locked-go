package controller

import (
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
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	if login.Email == "" || login.Password == "" {
		err = model.NewBadRequestApiError("invalid fields")

		http.Error(w, err.Message, err.StatusCode)

		return
	}
	user, err := c.Authenticate(login)

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	cookie := createEmailCookie(user.Email)

	http.SetCookie(w, cookie)

	redirect(w, req, "/")
}
