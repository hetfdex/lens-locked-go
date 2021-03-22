package main

import (
	"github.com/gorilla/mux"
	"lens-locked-go/controller"
	"net/http"
)

func main() {
	homeController := controller.NewHomeController()
	registerController := controller.NewRegisterController()

	r := mux.NewRouter()

	r.HandleFunc(homeController.Route, homeController.Handle)
	r.HandleFunc(registerController.Route, registerController.Handle).Methods(http.MethodGet)
	r.HandleFunc(registerController.Route, registerController.Register).Methods(http.MethodPost)

	err := http.ListenAndServe("localhost:8080", r)

	if err != nil {
		panic(err)
	}
}
