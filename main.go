package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"lens-locked-go/config"
	"lens-locked-go/controller"
	"lens-locked-go/middleware"
	"lens-locked-go/model"
	"lens-locked-go/rand"
	"lens-locked-go/repository"
	"lens-locked-go/service"
	"net/http"
)

func main() {
	cfg := config.LoadConfig(false)

	db := openDb(cfg)

	//resetDatabase(db)

	rp := repository.NewRepositories(db)

	sv := service.NewServices(rp, cfg)

	rt := mux.NewRouter()

	configureRouter(rt, cfg, sv)

	listenAndServe(rt, cfg.Server)
}

func configureRouter(rt *mux.Router, cfg *config.Config, sv *service.Services) {
	homeController := controller.NewHomeController()
	userController := controller.NewUserController(sv.User)
	galleryController := controller.NewGalleryController(sv.Gallery, sv.Image)
	dropboxController := controller.NewDropboxController(cfg.Dropbox, sv.Dropbox)

	authKey, err := rand.GenerateAuthKey()

	if err != nil {
		panic(err.Message)
	}
	csrfMdw := csrf.Protect(authKey, csrf.Secure(!cfg.Server.Debug))

	mdw := middleware.NewMiddleware(sv.User, userController.LoginRoute())

	rt.Use(csrfMdw)
	rt.Use(mdw.SetUser)

	rt.HandleFunc(homeController.HomeRoute(), homeController.HomeGet).Methods(http.MethodGet)

	rt.HandleFunc(userController.RegisterRoute(), userController.RegisterGet).Methods(http.MethodGet)
	rt.HandleFunc(userController.RegisterRoute(), userController.RegisterPost).Methods(http.MethodPost)

	rt.HandleFunc(userController.LoginRoute(), userController.LoginGet).Methods(http.MethodGet)
	rt.HandleFunc(userController.LoginRoute(), userController.LoginPost).Methods(http.MethodPost)

	rt.HandleFunc(userController.LogoutRoute(), userController.LogoutGet).Methods(http.MethodGet)

	rt.HandleFunc(galleryController.IndexRoute(), mdw.RequireUser(galleryController.IndexGet)).Methods(http.MethodGet)

	rt.HandleFunc(galleryController.CreateRoute(), mdw.RequireUser(galleryController.CreateGet)).Methods(http.MethodGet)
	rt.HandleFunc(galleryController.CreateRoute(), mdw.RequireUser(galleryController.CreatePost)).Methods(http.MethodPost)

	rt.HandleFunc(galleryController.EditRoute(), mdw.RequireUser(galleryController.EditGet)).Methods(http.MethodGet)
	rt.HandleFunc(galleryController.EditRoute(), mdw.RequireUser(galleryController.EditPost)).Methods(http.MethodPost)

	rt.HandleFunc(galleryController.UploadRoute(), mdw.RequireUser(galleryController.UploadGet)).Methods(http.MethodGet)
	rt.HandleFunc(galleryController.UploadRoute(), mdw.RequireUser(galleryController.UploadPost)).Methods(http.MethodPost)
	rt.HandleFunc(galleryController.UploadDropboxRoute(), mdw.RequireUser(galleryController.UploadDropboxPost)).Methods(http.MethodPost)

	rt.HandleFunc(galleryController.DeleteRoute(), mdw.RequireUser(galleryController.DeleteGet)).Methods(http.MethodGet)

	rt.HandleFunc(galleryController.GalleryRoute(), galleryController.GalleryGet).Methods(http.MethodGet)

	rt.HandleFunc(dropboxController.ConnectRoute(), mdw.RequireUser(dropboxController.ConnectGet)).Methods(http.MethodGet)
	rt.HandleFunc(dropboxController.CallbackRoute(), mdw.RequireUser(dropboxController.CallbackGet)).Methods(http.MethodGet)
	rt.HandleFunc(dropboxController.QueryRoute(), mdw.RequireUser(dropboxController.QueryGet)).Methods(http.MethodGet)

}

func openDb(cfg *config.Config) *gorm.DB {
	logLevel := logger.Warn

	if cfg.Server.Debug {
		logLevel = logger.Info
	}
	db, err := gorm.Open(cfg.Db.Dialector(), &gorm.Config{
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
	_ = db.Migrator().DropTable(&model.Image{})
	_ = db.Migrator().DropTable(&model.Dropbox{})

	_ = db.Migrator().CreateTable(&model.User{})
	_ = db.Migrator().CreateTable(&model.Gallery{})
	_ = db.Migrator().CreateTable(&model.Image{})
	_ = db.Migrator().CreateTable(&model.Dropbox{})
}

func listenAndServe(rt *mux.Router, sc *config.ServerConfig) {
	err := http.ListenAndServe(sc.Address, rt)

	if err != nil {
		panic(err)
	}
}
