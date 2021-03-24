package controller

import (
	"fmt"
	"lens-locked-go/config"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

type controller struct {
	Route       string
	view        *view.View
	userService *service.UserService
}

func (c *controller) Get(w http.ResponseWriter, req *http.Request) {
	apiErr := c.view.Render(w, nil)

	if apiErr != nil {
		http.Error(w, apiErr.Message, apiErr.StatusCode)
	}
	cookie, err := req.Cookie(config.CookieName)

	if err != nil {
		fmt.Println("no cookie available")

		return
	}

	user := &model.User{
		Token: cookie.Value,
	}

	user, apiErr = c.userService.AuthenticateWithToken(user)

	if apiErr != nil {
		http.Error(w, apiErr.Message, apiErr.StatusCode)

		return
	}
	fmt.Println(user.Email, "logged in")
}

func newController(route string, filename string, us *service.UserService) *controller {
	return &controller{
		Route:       route,
		view:        view.New(filename),
		userService: us,
	}
}
