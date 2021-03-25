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
	t, err := template.ParseFiles("view/base.gohtml", filename)

	if err != nil {
		panic(err)
	}
	return &View{
		template: t,
	}
}

func (v *View) Render(w http.ResponseWriter, data interface{}) *model.ApiError {
	w.Header().Set("Content-Type", "text/html")

	err := v.template.ExecuteTemplate(w, "base", data)

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}
