package controller

import (
	"lens-locked-go/model"
	"lens-locked-go/view"
	"net/http"
)

type homeController struct {
	homeView *view.View
}

func NewHomeController() *homeController {
	return &homeController{
		homeView: view.New(homeRoute, homeFilename),
	}
}

func (c *homeController) HomeGet(w http.ResponseWriter, req *http.Request) {
	data := &model.Data{}

	parseSuccessFromRoute(req, data)

	c.homeView.Render(w, req, data)
}

func (c *homeController) HomeRoute() string {
	return c.homeView.Route()
}
