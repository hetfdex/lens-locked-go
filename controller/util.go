package controller

import (
	"github.com/gorilla/schema"
	"net/http"
)

func parseForm(req *http.Request, result interface{}) error {
	err := req.ParseForm()

	if err != nil {
		return err
	}
	dec := schema.NewDecoder()

	err = dec.Decode(result, req.PostForm)

	if err != nil {
		return err
	}
	return nil
}
