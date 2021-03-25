package controller

import (
	"github.com/gorilla/schema"
	"lens-locked-go/config"
	"lens-locked-go/model"
	"lens-locked-go/validator"
	"net/http"
)

func parseForm(req *http.Request, result interface{}) *model.ApiError {
	err := req.ParseForm()

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	dec := schema.NewDecoder()

	err = dec.Decode(result, req.PostForm)

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func makeCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     config.CookieName,
		Value:    value,
		HttpOnly: true,
	}
}

func redirect(w http.ResponseWriter, req *http.Request, route string) {
	http.Redirect(w, req, route, http.StatusFound)
}

func validLoginForm(login *model.LoginForm) *model.ApiError {
	apiErr := validator.StringNotEmpty("email", login.Email)

	if apiErr != nil {
		return apiErr
	}
	apiErr = validator.StringNotEmpty("password", login.Password)

	if apiErr != nil {
		return apiErr
	}
	return nil
}

func validRegisterForm(register *model.RegisterForm) *model.ApiError {
	apiErr := validator.StringNotEmpty("name", register.Name)

	if apiErr != nil {
		return apiErr
	}
	apiErr = validator.StringNotEmpty("email", register.Email)

	if apiErr != nil {
		return apiErr
	}
	apiErr = validator.StringNotEmpty("password", register.Password)

	if apiErr != nil {
		return apiErr
	}
	return nil
}
