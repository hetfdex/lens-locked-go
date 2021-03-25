package controller

import (
	"github.com/gorilla/schema"
	"lens-locked-go/model"
	"lens-locked-go/util"
	"net/http"
)

func parseForm(req *http.Request, result interface{}) *model.ApiError {
	err := req.ParseForm()

	if err != nil {
		return model.NewBadRequestApiError(err.Error())
	}
	dec := schema.NewDecoder()

	err = dec.Decode(result, req.PostForm)

	if err != nil {
		return model.NewBadRequestApiError(err.Error())
	}
	return nil
}

func makeCookie(cookieValue string) (*http.Cookie, *model.ApiError) {
	if cookieValue == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("cookieValue"))
	}
	return &http.Cookie{
		Name:     util.CookieName,
		Value:    cookieValue,
		HttpOnly: true,
	}, nil
}

func redirect(w http.ResponseWriter, req *http.Request, route string) {
	if route == "" {
		return
	}
	http.Redirect(w, req, route, http.StatusFound)
}
