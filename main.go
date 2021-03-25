package main

import (
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lens-locked-go/config"
	"lens-locked-go/controller"
	"lens-locked-go/hash"
	"lens-locked-go/model"
	"lens-locked-go/repository"
	"lens-locked-go/service"
	"net/http"
)

func main() {
	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})

	resetDatabase(db)

	if err != nil {
		panic(err)
	}
	ur, apiErr := repository.NewUserRepository(db)

	hs := hash.New(config.HasherKey)

	us, apiErr := service.NewUserService(ur, hs)

	if apiErr != nil {
		panic(apiErr)
	}

	r := mux.NewRouter()

	configureRouter(us, r)

	err = http.ListenAndServe("localhost:8080", r)

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

func resetDatabase(db *gorm.DB) {
	dropUserTable(db)
	createUserTable(db)
}

func dropUserTable(db *gorm.DB) {
	_ = db.Migrator().DropTable(&model.User{})
}

func createUserTable(db *gorm.DB) {
	_ = db.Migrator().CreateTable(&model.User{})
}
