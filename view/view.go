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
	template *template.Template
}

func New(filename string) *View {
	if filename == "" {
		panic(errors.New(model.MustNotBeEmptyErrorMessage("filename")))
	}
	t, err := template.ParseFiles(baseFilename, alertFilename, filename)

	if err != nil {
		panic(err)
	}
	return &View{
		template: t,
	}
}

func (v *View) Render(w http.ResponseWriter, data *model.Data) *model.ApiError {
	w.Header().Set(contentTypeKey, contentTypeValue)

	err := v.template.ExecuteTemplate(w, baseTag, data)

	if err != nil {
		apiErr := model.NewInternalServerApiError(err.Error())

		alert := apiErr.Alert()

		http.Error(w, alert.Message, apiErr.StatusCode)
	}
	return nil
}
