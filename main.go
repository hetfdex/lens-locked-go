package main

import (
	"github.com/gorilla/mux"
	"lens-locked-go/config"
	"lens-locked-go/controller"
	"lens-locked-go/hash"
	"lens-locked-go/repository"
	"lens-locked-go/service"
	"net/http"
)

func main() {
	ur, apiErr := repository.NewUserRepository(config.Dsn)

	resetDatabase(ur)

	hs := hash.New(config.HasherKey)

	us, apiErr := service.NewUserService(ur, hs)

	if apiErr != nil {
		panic(apiErr)
	}

	r := mux.NewRouter()

	configureRouter(us, r)

	err := http.ListenAndServe("localhost:8080", r)

	if err != nil {
		panic(err)
	}
}

func configureRouter(us *service.UserService, r *mux.Router) {
	homeController := controller.NewHomeController(us)
	registerController := controller.NewRegisterController(us)
	loginController := controller.NewLoginController(us)

	r.HandleFunc(homeController.Route, homeController.Get).Methods(http.MethodGet)
	r.HandleFunc(registerController.Route, registerController.Get).Methods(http.MethodGet)
	r.HandleFunc(registerController.Route, registerController.Post).Methods(http.MethodPost)
	r.HandleFunc(loginController.Route, loginController.Get).Methods(http.MethodGet)
	r.HandleFunc(loginController.Route, loginController.Post).Methods(http.MethodPost)
}

func resetDatabase(ur *repository.UserRepository) {
	_ = ur.DropTable()
	_ = ur.CreateTable()
}
