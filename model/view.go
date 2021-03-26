package model

import (
	"errors"
	"html/template"
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

func NewView(filename string) *View {
	if filename == "" {
		panic(errors.New(MustNotBeEmptyErrorMessage("filename")))
	}
	t, err := template.ParseFiles(baseFilename, alertFilename, filename)

	if err != nil {
		panic(err)
	}
	return &View{
		template: t,
	}
}

func (v *View) Render(w http.ResponseWriter, data *Data) *ApiError {
	w.Header().Set(contentTypeKey, contentTypeValue)

	err := v.template.ExecuteTemplate(w, baseTag, data)

	if err != nil {
		return NewInternalServerApiError(err.Error())
	}
	return nil
}
