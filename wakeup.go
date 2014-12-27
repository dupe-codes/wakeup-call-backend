package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/njdup/wakeup-call-backend/conf"
	"github.com/njdup/wakeup-call-backend/models/user"

	"./controllers/userController"
)

func HomeHandler(res http.ResponseWriter, req *http.Request) {
	user := &user.User{
		Fullname: "Mr Dude",
	}
	err := user.Save()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(res, "This is a test! And the user was correctly saved hooray!")
}

// ConfigureRoutes sets all API routes
func configureRoutes(router *mux.Router) {
	router.HandleFunc("/test/", HomeHandler)

	router.HandleFunc("/users", userController.CreateUser).Methods("POST")
	router.HandleFunc("/users/login", userController.Login).Methods("POST")
}

// Main launches the API server
func main() {
	router := mux.NewRouter()
	configureRoutes(router)

	http.Handle("/", router)
	http.ListenAndServe(config.Settings.Port, nil)
}
