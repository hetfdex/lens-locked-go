package model

import (
	"errors"
	"html/template"
	"lens-locked-go/util"
	"net/http"
)

type Alert struct {
	Level   string
	Message string
}

type Data struct {
	Alert *Alert
	Data  interface{}
}

type View struct {
	template *template.Template
}

func NewView(filename string) *View {
	if filename == "" {
		panic(errors.New(util.MustNotBeEmptyErrorMessage("filename")))
	}
	t, err := template.ParseFiles(util.BaseFilename, util.AlertFilename, filename)

	if err != nil {
		panic(err)
	}
	return &View{
		template: t,
	}
}

func (v *View) Render(w http.ResponseWriter, data *Data) *ApiError {
	w.Header().Set(util.ContentTypeKey, util.ContentTypeValue)

	err := v.template.ExecuteTemplate(w, util.BaseTag, data)

	if err != nil {
		return NewInternalServerApiError(err.Error())
	}
	return nil
}
