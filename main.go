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
	cleanDatabase(us)

	r := mux.NewRouter()

	configureRoutes(us, r)

	err = http.ListenAndServe("localhost:8080", r)

	if err != nil {
		panic(err)
	}
}

func configureRoutes(us *service.UserService, r *mux.Router) {
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
