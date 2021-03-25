package main

import (
	"github.com/gorilla/mux"
	"lens-locked-go/repository"
	"lens-locked-go/service"
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
