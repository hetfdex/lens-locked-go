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

	db := openDb(cfg.Server, cfg.Postgres)

	//resetDatabase(db)

	ur := repository.NewUserRepository(db)
	ir := repository.NewImageRepository(db)
	gr := repository.NewGalleryRepository(db)

	us := service.NewUserService(ur, cfg.Crypto)
	is := service.NewImageService(ir)
	gs := service.NewGalleryService(gr)

	r := mux.NewRouter()

	configureRouter(r, cfg.Server, us, gs, is)

	listenAndServe(r, cfg.Server)
}

func openDb(sc *config.ServerConfig, dbc *config.PostgresConfig) *gorm.DB {
	logLevel := logger.Warn

	if sc.IsDebug {
		logLevel = logger.Info
	}
	db, err := gorm.Open(dbc.Dialector(), &gorm.Config{
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
	_ = db.Migrator().CreateTable(&model.User{})
	_ = db.Migrator().CreateTable(&model.Gallery{})
	_ = db.Migrator().CreateTable(&model.Image{})
}

func configureRouter(r *mux.Router, sc *config.ServerConfig, us service.IUserService, gs service.IGalleryService, is service.IImageService) {
	homeController := controller.NewHomeController()
	userController := controller.NewUserController(us)
	galleryController := controller.NewGalleryController(gs, is)

	authKey, err := rand.GenerateAuthKey()

	if err != nil {
		panic(err.Message)
	}
	csrfMdw := csrf.Protect(authKey, csrf.Secure(sc.IsDebug))

	mdw := middleware.NewMiddleware(us, userController.LoginRoute())

	r.Use(csrfMdw)
	r.Use(mdw.SetUser)

	r.HandleFunc(homeController.HomeRoute(), homeController.HomeGet).Methods(http.MethodGet)
	r.HandleFunc(userController.RegisterRoute(), userController.RegisterGet).Methods(http.MethodGet)
	r.HandleFunc(userController.RegisterRoute(), userController.RegisterPost).Methods(http.MethodPost)
	r.HandleFunc(userController.LoginRoute(), userController.LoginGet).Methods(http.MethodGet)
	r.HandleFunc(userController.LoginRoute(), userController.LoginPost).Methods(http.MethodPost)
	r.HandleFunc(userController.LogoutRoute(), userController.LogoutGet).Methods(http.MethodGet)
	r.HandleFunc(galleryController.IndexRoute(), mdw.RequireUser(galleryController.IndexGet)).Methods(http.MethodGet)
	r.HandleFunc(galleryController.CreateRoute(), mdw.RequireUser(galleryController.CreateGet)).Methods(http.MethodGet)
	r.HandleFunc(galleryController.CreateRoute(), mdw.RequireUser(galleryController.CreatePost)).Methods(http.MethodPost)
	r.HandleFunc(galleryController.EditRoute(), mdw.RequireUser(galleryController.EditGet)).Methods(http.MethodGet)
	r.HandleFunc(galleryController.EditRoute(), mdw.RequireUser(galleryController.EditPost)).Methods(http.MethodPost)
	r.HandleFunc(galleryController.UploadRoute(), mdw.RequireUser(galleryController.UploadGet)).Methods(http.MethodGet)
	r.HandleFunc(galleryController.UploadRoute(), mdw.RequireUser(galleryController.UploadPost)).Methods(http.MethodPost)
	r.HandleFunc(galleryController.DeleteRoute(), mdw.RequireUser(galleryController.DeleteGet)).Methods(http.MethodGet)
	r.HandleFunc(galleryController.GalleryRoute(), galleryController.GalleryGet).Methods(http.MethodGet)
}

func listenAndServe(r *mux.Router, sc *config.ServerConfig) {
	err := http.ListenAndServe(sc.Address, r)

	if err != nil {
		panic(err)
	}
}
