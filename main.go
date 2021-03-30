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

	r.HandleFunc(homeController.HomeRoute(), homeController.GetHome).Methods(http.MethodGet)
	r.HandleFunc(userController.RegisterUserRoute(), userController.GetRegisterUser).Methods(http.MethodGet)
	r.HandleFunc(userController.RegisterUserRoute(), userController.PostRegisterUser).Methods(http.MethodPost)
	r.HandleFunc(userController.LoginUserRoute(), userController.GetLoginUser).Methods(http.MethodGet)
	r.HandleFunc(userController.LoginUserRoute(), userController.PostLoginUser).Methods(http.MethodPost)
	r.HandleFunc(galleryController.IndexGalleryRoute(), mdw.RequireUser(galleryController.GetIndexGallery)).Methods(http.MethodGet)
	r.HandleFunc(galleryController.CreateGalleryRoute(), mdw.RequireUser(galleryController.GetCreateGallery)).Methods(http.MethodGet)
	r.HandleFunc(galleryController.CreateGalleryRoute(), mdw.RequireUser(galleryController.PostCreateGallery)).Methods(http.MethodPost)
	r.HandleFunc(galleryController.EditGalleryRoute(), mdw.RequireUser(galleryController.GetEditGallery)).Methods(http.MethodGet)
	r.HandleFunc(galleryController.EditGalleryRoute(), mdw.RequireUser(galleryController.PostEditGallery)).Methods(http.MethodPost)
	r.HandleFunc(galleryController.DeleteGalleryRoute(), mdw.RequireUser(galleryController.GetDeleteGallery)).Methods(http.MethodGet)
	r.HandleFunc(galleryController.GalleryRoute(), galleryController.GetGallery).Methods(http.MethodGet).Name(controller.GalleryRouteName)
}

func listenAndServe(r *mux.Router) {
	err := http.ListenAndServe(address, r)

	if err != nil {
		panic(err)
	}
}
