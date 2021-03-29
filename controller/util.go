package controller

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"lens-locked-go/context"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"lens-locked-go/view"
	"net/http"
)

const CookieName = "login_token"

func Redirect(w http.ResponseWriter, req *http.Request, route string) {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	http.Redirect(w, req, route, http.StatusFound)
}

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

func handleError(w http.ResponseWriter, view *view.View, err *model.Error, data *model.DataView) {
	data.Alert = err.Alert()

	w.WriteHeader(err.StatusCode)

	view.Render(w, data)
}

func makeUrl(router *mux.Router, routeName string, key string, value string) (string, *model.Error) {
	url, err := router.Get(routeName).URL(key, value)

	if err != nil {
		er := model.NewInternalServerApiError(err.Error())

		return "", er
	}
	return url.String(), nil
}

func getGallery(req *http.Request, gs service.IGalleryService) (*model.Gallery, *model.Error) {
	vars := mux.Vars(req)

	id := uuid.FromStringOrNil(vars[idKey])

	if id == uuid.Nil {
		err := model.NewBadRequestApiError(invalidUUIDErrorMessage)

		return nil, err
	}
	gallery, err := gs.Get(id)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func getGalleryWithPermission(req *http.Request, gs service.IGalleryService) (*model.Gallery, *model.Error) {
	gallery, err := getGallery(req, gs)

	if err != nil {
		return nil, err
	}
	user, err := context.User(req.Context())

	if err != nil {
		return nil, err
	}

	if gallery.UserId != user.ID {
		err = model.NewForbiddenApiError(userNotOwnerErrorMessage)

		return nil, err
	}
	return gallery, nil
}
