package view

import (
	"errors"
	"html/template"
	"lens-locked-go/context"
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

func (v *View) Render(w http.ResponseWriter, req *http.Request, data *model.Data) {
	w.Header().Set(contentTypeKey, contentTypeValue)

	user, _ := context.User(req.Context())

	if user != nil {
		data.User = user
	}
	err := v.template.ExecuteTemplate(w, baseTag, data)

	if err != nil {
		er := model.NewInternalServerApiError(err.Error())

		http.Error(w, er.Message, er.StatusCode)
	}
}

func (v *View) Route() string {
	return v.route
}
