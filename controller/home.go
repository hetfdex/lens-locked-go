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
	data := &model.Data{}

	parseSuccessFromRoute(req, data)

	c.homeView.Render(w, req, data)
}

func HomeRoute() string {
	return homeRoute
}
