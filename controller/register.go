package controller

import (
	"fmt"
	"lens-locked-go/model"
	"net/http"
)

type registerController struct {
	*controller
}

func NewRegisterController() *registerController {
	return &registerController{
		newController("/register", "view/register.gohtml"),
	}
}

func (c *registerController) Register(w http.ResponseWriter, req *http.Request) {
	form := &model.Register{}

	err := parseForm(req, form)

	if err != nil {
		panic(err)
	}
	_, _ = fmt.Fprintln(w, *form)
}
