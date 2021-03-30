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

func (c *homeController) GetHome(w http.ResponseWriter, req *http.Request) {
	viewData := &model.DataView{}

	parseSuccessRoute(req, viewData)

	c.homeView.Render(w, viewData)
}

func (c *homeController) HomeRoute() string {
	return c.homeView.Route()
}
