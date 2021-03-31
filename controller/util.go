package controller

import (
	"github.com/gorilla/schema"
	"lens-locked-go/model"
	"lens-locked-go/view"
	"net/http"
)

const successRouteQuery = "?success="

const successRouteKey = "success"

const registerUserValue = "register"
const loginUserValue = "login"
const logoutUserValue = "logout"
const createGalleryValue = "createGallery"
const editGalleryValue = "editGallery"
const uploadGalleryValue = "uploadGallery"
const deleteGalleryValue = "deleteGallery"

const registerUserSuccessMessage = "Registration Successful"
const loginUserSuccessMessage = "Login Successful"
const logoutUserSuccessMessage = "Logout Successful"
const createGallerySuccessMessage = "Created Successfully"
const editGallerySuccessMessage = "Edited Successfully"
const uploadGallerySuccessMessage = "Uploaded Successfully"
const deleteGallerySuccessMessage = "Deleted Successfully"

func parseForm(req *http.Request, result interface{}) *model.Error {
	err := req.ParseForm()

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	dec := schema.NewDecoder()

	dec.IgnoreUnknownKeys(true)

	err = dec.Decode(result, req.PostForm)

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func addSuccessToRoute(route string, value string) string {
	return route + successRouteQuery + value
}

func parseSuccessFromRoute(req *http.Request, data *model.Data) {
	param := req.URL.Query().Get(successRouteKey)

	switch param {
	case registerUserValue:
		data.Alert = model.NewSuccessAlert(registerUserSuccessMessage)
	case loginUserValue:
		data.Alert = model.NewSuccessAlert(loginUserSuccessMessage)
	case logoutUserValue:
		data.Alert = model.NewSuccessAlert(logoutUserSuccessMessage)
	case createGalleryValue:
		data.Alert = model.NewSuccessAlert(createGallerySuccessMessage)
	case editGalleryValue:
		data.Alert = model.NewSuccessAlert(editGallerySuccessMessage)
	case uploadGalleryValue:
		data.Alert = model.NewSuccessAlert(uploadGallerySuccessMessage)
	case deleteGalleryValue:
		data.Alert = model.NewSuccessAlert(deleteGallerySuccessMessage)
	}
}

func handleError(w http.ResponseWriter, req *http.Request, data *model.Data, err *model.Error, view *view.View) {
	data.Alert = err.Alert()

	w.WriteHeader(err.StatusCode)

	view.Render(w, req, data)
}
