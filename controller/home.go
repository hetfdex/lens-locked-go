package controller

import (
	"lens-locked-go/model"
	"lens-locked-go/view"
	"net/http"
)

const homeRoute = "/"
const homeFilename = "view/home.gohtml"

type homeController struct {
	homeView *view.View
}

func NewHomeController() *homeController {
	return &homeController{
		homeView: view.New(homeRoute, homeFilename),
	}
}

func (c *homeController) HomeGet(w http.ResponseWriter, _ *http.Request) {
	data := &model.DataView{}

	c.homeView.Render(w, data)
}

func (c *homeController) HomeRoute() string {
	return c.homeView.Route()
}
