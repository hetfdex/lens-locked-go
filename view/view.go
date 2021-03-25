package view

import (
	"errors"
	"html/template"
	"lens-locked-go/model"
	"lens-locked-go/util"
	"net/http"
)

type View struct {
	template *template.Template
}

func New(filename string) *View {
	if filename == "" {
		panic(errors.New(util.MustNotBeEmptyErrorMessage("filename")))
	}
	t, err := template.ParseFiles(util.BaseFilename, filename)

	if err != nil {
		panic(err)
	}
	return &View{
		template: t,
	}
}

func (v *View) Render(w http.ResponseWriter, data interface{}) *model.ApiError {
	w.Header().Set(util.ContentTypeKey, util.ContentTypeValue)

	err := v.template.ExecuteTemplate(w, util.BaseTag, data)

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}
