package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"lens-locked-go/controller"
	"lens-locked-go/middleware"
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

	//resetDatabase(db)

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
	homeController := controller.NewHomeController()
	userController := controller.NewUserController(us)
	galleryController := controller.NewGalleryController(r, gs)

	mdw := middleware.NewMiddleware(us)

	r.HandleFunc(homeController.HomeRoute(), homeController.HomeGet).Methods(http.MethodGet)
	r.HandleFunc(userController.RegisterRoute(), userController.RegisterGet).Methods(http.MethodGet)
	r.HandleFunc(userController.RegisterRoute(), userController.RegisterPost).Methods(http.MethodPost)
	r.HandleFunc(userController.LoginRoute(), userController.LoginGet).Methods(http.MethodGet)
	r.HandleFunc(userController.LoginRoute(), userController.LoginPost).Methods(http.MethodPost)
	r.HandleFunc(galleryController.GalleryRoute(), galleryController.GalleryGet).Methods(http.MethodGet).Name(controller.GalleryRouteName)
	r.HandleFunc(galleryController.CreateRoute(), mdw.RequireUser(galleryController.CreateGet)).Methods(http.MethodGet)
	r.HandleFunc(galleryController.CreateRoute(), mdw.RequireUser(galleryController.CreatePost)).Methods(http.MethodPost)
	r.HandleFunc(galleryController.EditRoute(), mdw.RequireUser(galleryController.EditGet)).Methods(http.MethodGet)
	r.HandleFunc(galleryController.EditRoute(), mdw.RequireUser(galleryController.EditPost)).Methods(http.MethodPost)
}

func listenAndServe(r *mux.Router) {
	err := http.ListenAndServe(address, r)

	if err != nil {
		panic(err)
	}
}
