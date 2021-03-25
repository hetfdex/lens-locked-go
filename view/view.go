package view

import (
	"html/template"
	"lens-locked-go/model"
	"lens-locked-go/validator"
	"net/http"
)

type View struct {
	template *template.Template
}

func New(filename string) (*View, *model.ApiError) {
	apiErr := validator.StringNotEmpty("filename", filename)

	if apiErr != nil {
		return nil, apiErr
	}
	t, err := template.ParseFiles("view/base.gohtml", filename)

	if err != nil {
		panic(err)
	}
	return &View{
		template: t,
	}, nil
}

func (v *View) Render(w http.ResponseWriter, data interface{}) *model.ApiError {
	w.Header().Set("Content-Type", "text/html")

	err := v.template.ExecuteTemplate(w, "base", data)

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}
