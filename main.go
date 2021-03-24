package main

import (
	"github.com/gorilla/mux"
	"lens-locked-go/config"
	"lens-locked-go/controller"
	"lens-locked-go/service"
	"net/http"
)

func main() {
	us, apiErr := service.New(config.Dsn)

	if apiErr != nil {
		panic(apiErr)
	}
	cleanDatabase(us)

	r := mux.NewRouter()

	configureRouter(us, r)

	err := http.ListenAndServe("localhost:8080", r)

	if err != nil {
		panic(err)
	}
}

func configureRouter(us *service.UserService, r *mux.Router) {
	homeController := controller.NewHomeController()
	registerController := controller.NewRegisterController(us)
	loginController := controller.NewLoginController(us)

	r.HandleFunc(homeController.Route, homeController.Handle).Methods(http.MethodGet)
	r.HandleFunc(registerController.Route, registerController.Handle).Methods(http.MethodGet)
	r.HandleFunc(registerController.Route, registerController.Register).Methods(http.MethodPost)
	r.HandleFunc(loginController.Route, loginController.Handle).Methods(http.MethodGet)
	r.HandleFunc(loginController.Route, loginController.Login).Methods(http.MethodPost)
}

func cleanDatabase(us *service.UserService) {
	_ = us.DropTable()
	_ = us.CreateTable()
}
