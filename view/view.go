package view

import (
	"html/template"
	"lens-locked-go/model"
	"net/http"
)

type View struct {
	*template.Template
}

func New(filename string) *View {
	t, err := template.ParseFiles("view/base.gohtml", filename)

	if err != nil {
		panic(err)
	}
	return &View{
		t,
	}
}

func (v *View) Render(w http.ResponseWriter, data interface{}) *model.ApiError {
	w.Header().Set("Content-Type", "text/html")

	err := v.ExecuteTemplate(w, "base", data)

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}
