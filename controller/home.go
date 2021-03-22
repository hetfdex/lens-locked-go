package controller

type homeController struct {
	*controller
}

func NewHomeController() *homeController {
	return &homeController{
		newController("/", "view/home.gohtml"),
	}
}
