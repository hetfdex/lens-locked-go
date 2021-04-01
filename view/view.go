package view

import (
	"errors"
	"github.com/gorilla/csrf"
	"html/template"
	"lens-locked-go/context"
	"lens-locked-go/model"
	"lens-locked-go/util"
	"net/http"
)

type View struct {
	route    string
	template *template.Template
}

func New(route string, filename string) *View {
	if route == "" {
		panic(errors.New(util.MustNotBeEmptyErrorMessage("route")))
	}
	if filename == "" {
		panic(errors.New(util.MustNotBeEmptyErrorMessage("filename")))
	}
	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return "<h1>csrfField</h1>"
		},
	}).ParseFiles("view/base.gohtml", "view/navbar.gohtml", "view/alert.gohtml", filename)

	if err != nil {
		panic(err)
	}
	return &View{
		route:    route,
		template: t,
	}
}

func (v *View) Render(w http.ResponseWriter, req *http.Request, data *model.Data) {
	w.Header().Set("Content-Type", "text/html")

	user, _ := context.User(req.Context())

	if user != nil {
		data.User = user
	}
	csrfField := csrf.TemplateField(req)

	tpl := v.template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})

	err := tpl.ExecuteTemplate(w, "base", data)

	if err != nil {
		er := model.NewInternalServerApiError(err.Error())

		http.Error(w, er.Message, er.StatusCode)
	}
}

func (v *View) Route() string {
	return v.route
}
