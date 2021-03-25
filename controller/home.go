package controller

import (
	"lens-locked-go/service"
)

type homeController struct {
	*controller
}

func NewHomeController(us service.IUserService) *homeController {
	return &homeController{
		newController("/", "view/home.gohtml", us),
	}
}
