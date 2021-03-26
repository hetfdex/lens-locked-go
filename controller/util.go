package controller

import (
	"github.com/gorilla/schema"
	"lens-locked-go/model"
	"lens-locked-go/view"
	"net/http"
)

const name = "login_token"

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
		Name:     name,
		Value:    value,
		HttpOnly: true,
	}, nil
}

func redirect(w http.ResponseWriter, req *http.Request, route string) {
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
