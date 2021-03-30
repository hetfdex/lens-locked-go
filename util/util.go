package util

import (
	"errors"
	"lens-locked-go/model"
	"net/http"
)

const CookieName = "login_token"

func Redirect(w http.ResponseWriter, req *http.Request, route string) {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	http.Redirect(w, req, route, http.StatusFound)
}
