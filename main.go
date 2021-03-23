package main

import (
	"github.com/gorilla/mux"
	"lens-locked-go/config"
	"lens-locked-go/controller"
	"lens-locked-go/service"
	"net/http"
)

func main() {
	us, err := service.New(config.Dsn)

	if err != nil {
		panic(err)
	}
	_ = us.DropTable()
	_ = us.CreateTable()

	homeController := controller.NewHomeController()
	registerController := controller.NewRegisterController(us)

	r := mux.NewRouter()

	r.HandleFunc(homeController.Route, homeController.Handle)
	r.HandleFunc(registerController.Route, registerController.Handle).Methods(http.MethodGet)
	r.HandleFunc(registerController.Route, registerController.Register).Methods(http.MethodPost)

	err = http.ListenAndServe("localhost:8080", r)

	if err != nil {
		panic(err)
	}
}
