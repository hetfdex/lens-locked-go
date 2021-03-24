package controller

import (
	"github.com/gorilla/schema"
	"lens-locked-go/model"
	"net/http"
)

func parseForm(req *http.Request, result interface{}) *model.ApiError {
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

func createEmailCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:  "email",
		Value: value,
	}
}

func redirect(w http.ResponseWriter, req *http.Request, route string) {
	http.Redirect(w, req, route, http.StatusFound)
}
