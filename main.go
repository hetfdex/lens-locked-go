package main

import (
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lens-locked-go/config"
	"lens-locked-go/hash"
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
