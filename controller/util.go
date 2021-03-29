package controller

import (
	"github.com/gorilla/schema"
	"lens-locked-go/model"
	"lens-locked-go/view"
	"net/http"
)

const homeRoute = "/"
const homeFilename = "view/home.gohtml"

const LoginUserRoute = "/login"
const loginUserFilename = "view/user_login.gohtml"

const registerUserRoute = "/register"
const registerUserFilename = "view/user_register.gohtml"

const createGalleryRoute = "/gallery/create"
const createGalleryFilename = "view/gallery_create.gohtml"

const galleryRoute = "/gallery/:id"
const galleryFilename = "view/gallery.gohtml"

const noUserInContextErrorMessage = "no user found in context"

const CookieName = "login_token"

func parseForm(req *http.Request, result interface{}) *model.Error {
	err := req.ParseForm()

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	dec := schema.NewDecoder()

	err = dec.Decode(result, req.PostForm)

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func makeCookie(value string) (*http.Cookie, *model.Error) {
	if value == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("value"))
	}
	return &http.Cookie{
		Name:     CookieName,
		Value:    value,
		HttpOnly: true,
	}, nil
}

func Redirect(w http.ResponseWriter, req *http.Request, route string) {
	if route == "" {
		return
	}
	http.Redirect(w, req, route, http.StatusFound)
}

func handleError(view *view.View, w http.ResponseWriter, err *model.Error, data *model.DataView) {
	data.Alert = err.Alert()

	w.WriteHeader(err.StatusCode)

	view.Render(w, data)
}
