package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"lens-locked-go/config"
	"lens-locked-go/controller"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"net/http"
)

func main() {
	us, err := service.New(config.Dsn)

	if err != nil {
		panic(err)
	}
	err = us.DropTable()

	if err != nil {
		panic(err)
	}
	err = us.CreateTable()

	if err != nil {
		panic(err)
	}

	u := &model.User{
		Name:  "Jose",
		Email: "test@email.com",
	}

	err = us.Create(u)

	if err != nil {
		panic(err)
	}
	result, err := us.Read("id", u.ID)

	if err != nil {
		panic(err)
	}
	fmt.Println(result.ID, result.Name, result.Email)

	u.Email = "fake@email.com"

	err = us.Update(u)

	result, err = us.Read("email", u.Email)

	if err != nil {
		panic(err)
	}
	fmt.Println(result.ID, result.Name, result.Email)

	err = us.Delete(u.ID)

	if err != nil {
		panic(err)
	}
	_, _ = us.Read("id", u.ID)

	homeController := controller.NewHomeController()
	registerController := controller.NewRegisterController()

	r := mux.NewRouter()

	r.HandleFunc(homeController.Route, homeController.Handle)
	r.HandleFunc(registerController.Route, registerController.Handle).Methods(http.MethodGet)
	r.HandleFunc(registerController.Route, registerController.Register).Methods(http.MethodPost)

	err = http.ListenAndServe("localhost:8080", r)

	if err != nil {
		panic(err)
	}
}
