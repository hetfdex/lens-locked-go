package util

import (
	"errors"
	"fmt"
	"net/http"
)

const CookieName = "login_token"

func Redirect(w http.ResponseWriter, req *http.Request, route string) {
	if route == "" {
		panic(errors.New(MustNotBeEmptyErrorMessage("route")))
	}
	http.Redirect(w, req, route, http.StatusFound)
}

func MustNotBeEmptyErrorMessage(value string) string {
	message := fmt.Sprintf("%s must not be empty", value)

	return message
}

func InUseErrorMessage(value string) string {
	message := fmt.Sprintf("%s is already in use", value)

	return message
}

func InvalidErrorMessage(value string) string {
	message := fmt.Sprintf("invalid %s", value)

	return message
}
