package controller

import (
	"fmt"
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
	registration := &model.Register{}

	err := parseForm(req, registration)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	user := &model.User{
		Name:  registration.Name,
		Email: registration.Email,
	}

	err = c.Create(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, _ = fmt.Fprintln(w, user.Name)
	_, _ = fmt.Fprintln(w, user.Email)
}
