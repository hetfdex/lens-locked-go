package view

import (
	"html/template"
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

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")

	return v.ExecuteTemplate(w, "base", data)
}
