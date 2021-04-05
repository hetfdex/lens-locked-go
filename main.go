package main

import (
	"context"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
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
	"time"
)

func main() {
	cfg := config.LoadConfig(false)

	db := openDb(cfg)

	resetDatabase(db)

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
	//oAuthController := controller.NewOAuthController(sv.User, sv.OAuth)

	authKey, err := rand.GenerateAuthKey()

	if err != nil {
		panic(err.Message)
	}
	csrfMdw := csrf.Protect(authKey, csrf.Secure(!cfg.Server.Debug))

	mdw := middleware.NewMiddleware(sv.User, userController.LoginRoute())

	rt.Use(csrfMdw)
	rt.Use(mdw.SetUser)

	dbxOauth := &oauth2.Config{
		ClientID:     cfg.OAuth.Id,
		ClientSecret: cfg.OAuth.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.OAuth.AuthUrl,
			TokenURL: cfg.OAuth.TokenUrl,
		},
		RedirectURL: cfg.OAuth.RedirectUrl,
	}

	dbxConnect := func(w http.ResponseWriter, req *http.Request) {
		state := csrf.Token(req)

		cookie := &http.Cookie{
			Name:     "oauth_token",
			Value:    state,
			Expires:  time.Now().Add(time.Minute * 5),
			HttpOnly: true,
		}

		url := dbxOauth.AuthCodeURL(state)

		http.SetCookie(w, cookie)

		http.Redirect(w, req, url, http.StatusFound)
	}

	dbxCallback := func(w http.ResponseWriter, req *http.Request) {
		errr := req.ParseForm()

		if errr != nil {
			http.Error(w, errr.Error(), http.StatusFailedDependency)

			return
		}
		state := req.FormValue("state")

		cookie, errr := req.Cookie("oauth_token")

		if errr != nil {
			http.Error(w, errr.Error(), http.StatusFailedDependency)

			return
		}

		if state != cookie.Value {
			http.Error(w, "invalid state", http.StatusFailedDependency)

			return
		}
		cookie.Value = ""
		cookie.Expires = time.Now().Add(-time.Hour)

		http.SetCookie(w, cookie)

		code := req.FormValue("code")

		token, errr := dbxOauth.Exchange(context.TODO(), code)

		if errr != nil {
			http.Error(w, errr.Error(), http.StatusFailedDependency)
		}
		fmt.Printf("%+v", token)
	}

	connectPath := "/oauth/dropbox/connect"
	callbackPath := "/oauth/dropbox/callback"

	rt.HandleFunc(connectPath, mdw.RequireUser(dbxConnect)).Methods(http.MethodGet)
	rt.HandleFunc(callbackPath, mdw.RequireUser(dbxCallback)).Methods(http.MethodGet)

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

	rt.HandleFunc(galleryController.DeleteRoute(), mdw.RequireUser(galleryController.DeleteGet)).Methods(http.MethodGet)

	rt.HandleFunc(galleryController.GalleryRoute(), galleryController.GalleryGet).Methods(http.MethodGet)
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
	_ = db.Migrator().DropTable(&model.OAuth{})

	_ = db.Migrator().CreateTable(&model.User{})
	_ = db.Migrator().CreateTable(&model.Gallery{})
	_ = db.Migrator().CreateTable(&model.Image{})
	_ = db.Migrator().CreateTable(&model.OAuth{})
}

func listenAndServe(rt *mux.Router, sc *config.ServerConfig) {
	err := http.ListenAndServe(sc.Address, rt)

	if err != nil {
		panic(err)
	}
}
