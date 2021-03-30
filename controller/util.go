package controller

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"lens-locked-go/context"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"net/http"
)

const cookieName = "login_token"
const invalidTokenValue = "invalidTokenValue"

const successRouteQuery = "?success="
const successRouteKey = "success"
const registerUserValue = "register"
const loginUserValue = "login"
const logoutUserValue = "logout"
const createGalleryValue = "createGallery"
const editGalleryValue = "editGallery"
const deleteGalleryValue = "deleteGallery"
const registerUserSuccessMessage = "Registration Successful"
const loginUserSuccessMessage = "Login Successful"
const logoutUserSuccessMessage = "Logout Successful"
const createGallerySuccessMessage = "Created Successfully"
const editGallerySuccessMessage = "Edited Successfully"
const deleteGallerySuccessMessage = "Deleted Successfully"

func Redirect(w http.ResponseWriter, req *http.Request, route string) {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	http.Redirect(w, req, route, http.StatusFound)
}

func CookieName() string {
	return cookieName
}

func makeSuccessRoute(route string, value string) string {
	return route + successRouteQuery + value
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

func parseSuccessRoute(req *http.Request, viewData *model.DataView) {
	param := req.URL.Query().Get(successRouteKey)

	if param == registerUserValue {
		viewData.Alert = model.NewSuccessAlert(registerUserSuccessMessage)
	} else if param == loginUserValue {
		viewData.Alert = model.NewSuccessAlert(loginUserSuccessMessage)
	} else if param == logoutUserValue {
		viewData.Alert = model.NewSuccessAlert(logoutUserSuccessMessage)
	} else if param == createGalleryValue {
		viewData.Alert = model.NewSuccessAlert(createGallerySuccessMessage)
	} else if param == editGalleryValue {
		viewData.Alert = model.NewSuccessAlert(editGallerySuccessMessage)
	} else if param == deleteGalleryValue {
		viewData.Alert = model.NewSuccessAlert(deleteGallerySuccessMessage)
	}
}

func makeCookie(value string) (*http.Cookie, *model.Error) {
	if value == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("value"))
	}
	return &http.Cookie{
		Name:     cookieName,
		Value:    value,
		HttpOnly: true,
	}, nil
}

func getGallery(req *http.Request, gs service.IGalleryService) (*model.Gallery, *model.Error) {
	vars := mux.Vars(req)

	id := uuid.FromStringOrNil(vars[idKey])

	if id == uuid.Nil {
		err := model.NewBadRequestApiError(invalidUUIDErrorMessage)

		return nil, err
	}
	gallery, err := gs.GetById(id)

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
