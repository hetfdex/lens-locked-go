package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"lens-locked-go/controller"
	"lens-locked-go/model"
	"lens-locked-go/repository"
	"lens-locked-go/service"
	"net/http"
)

const host = "localhost"
const port = 5432
const user = "postgres"
const password = "Abcde12345!"
const dbname = "lenslocked_dev"

var dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

const address = "localhost:8080"

func main() {
	db := openDb(true)

	resetDatabase(db)

	ur := repository.NewUserRepository(db)
	gr := repository.NewGalleryRepository(db)

	us := service.NewUserService(ur)
	gs := service.NewGalleryService(gr)

	r := mux.NewRouter()

	configureRouter(r, us, gs)

	listenAndServe(r)
}

func openDb(isDebug bool) *gorm.DB {
	var logLevel = logger.Warn

	if isDebug {
		logLevel = logger.Info
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		panic(err)
	}
	return db
}

func resetDatabase(db *gorm.DB) {
	_ = db.Migrator().DropTable(&model.User{})
	_ = db.Migrator().DropTable(&model.Gallery{})
	_ = db.Migrator().CreateTable(&model.User{})
	_ = db.Migrator().CreateTable(&model.Gallery{})
}

func configureRouter(r *mux.Router, us service.IUserService, gs service.IGalleryService) {
	homeController := controller.NewHomeController(us)
	registerController := controller.NewRegisterController(us)
	loginController := controller.NewLoginController(us)
	galleryController := controller.NewGalleryController(gs)

	r.HandleFunc(homeController.Route, homeController.Get).Methods(http.MethodGet)
	r.HandleFunc(registerController.Route, registerController.Get).Methods(http.MethodGet)
	r.HandleFunc(registerController.Route, registerController.Post).Methods(http.MethodPost)
	r.HandleFunc(loginController.Route, loginController.Get).Methods(http.MethodGet)
	r.HandleFunc(loginController.Route, loginController.Post).Methods(http.MethodPost)
	r.HandleFunc(galleryController.Route, galleryController.Get).Methods(http.MethodGet)

}

func listenAndServe(r *mux.Router) {
	err := http.ListenAndServe(address, r)

	if err != nil {
		panic(err)
	}
}
