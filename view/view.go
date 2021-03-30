package view

import (
	"errors"
	"html/template"
	"lens-locked-go/model"
	"net/http"
)

const baseTag = "base"
const baseFilename = "view/base.gohtml"
const alertFilename = "view/alert.gohtml"

const contentTypeKey = "Content-Type"
const contentTypeValue = "text/html"

type View struct {
	route    string
	template *template.Template
}

func New(route string, filename string) *View {
	if route == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("route")))
	}
	if filename == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("filename")))
	}
	t, err := template.ParseFiles(baseFilename, alertFilename, filename)

	if err != nil {
		panic(err)
	}
	return &View{
		route:    route,
		template: t,
	}
}

func (v *View) Render(w http.ResponseWriter, viewData *model.DataView) {
	w.Header().Set(contentTypeKey, contentTypeValue)

	err := v.template.ExecuteTemplate(w, baseTag, viewData)

	if err != nil {
		er := model.NewInternalServerApiError(err.Error())

		http.Error(w, er.Message, er.StatusCode)
	}
}

func (v *View) Route() string {
	return v.route
}
