package main

import (
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lens-locked-go/controller"
	"lens-locked-go/model"
	"lens-locked-go/repository"
	"lens-locked-go/service"
	"lens-locked-go/util"
	"net/http"
)

func main() {
	db := openDb()

	resetDatabase(db)

	ur := repository.NewUserRepository(db)

	us := service.NewUserService(ur)

	r := mux.NewRouter()

	configureRouter(us, r)

	listenAndServe(r)
}

func openDb() *gorm.DB {
	db, err := gorm.Open(postgres.Open(util.Dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	return db
}

func resetDatabase(db *gorm.DB) {
	_ = db.Migrator().DropTable(&model.User{})
	_ = db.Migrator().CreateTable(&model.User{})
}

func configureRouter(us service.IUserService, r *mux.Router) {
	homeController := controller.NewHomeController(us)
	registerController := controller.NewRegisterController(us)
	loginController := controller.NewLoginController(us)

	r.HandleFunc(homeController.Route, homeController.Get).Methods(http.MethodGet)
	r.HandleFunc(registerController.Route, registerController.Get).Methods(http.MethodGet)
	r.HandleFunc(registerController.Route, registerController.Post).Methods(http.MethodPost)
	r.HandleFunc(loginController.Route, loginController.Get).Methods(http.MethodGet)
	r.HandleFunc(loginController.Route, loginController.Post).Methods(http.MethodPost)
}

func listenAndServe(r *mux.Router) {
	err := http.ListenAndServe(util.Address, r)

	if err != nil {
		panic(err)
	}
}
