package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/njdup/wakeup-call-backend/conf"
	"github.com/njdup/wakeup-call-backend/models/user"

	"./controllers/userController"
)

// ConfigureRoutes sets all API routes
func configureRoutes(router *mux.Router, sessionStore *sessions.CookieStore) {
	router.Handle("/users", userController.CreateUser(sessionStore)).Methods("POST")
	router.Handle("/users/login", userController.Login(sessionStore)).Methods("POST")
	router.Handle("/users/logout", userController.Logout(sessionStore)).Methods("POST")
	router.Handle("/users", userController.AllUsers(sessionStore)).Methods("GET")
}

// Main launches the API server
func main() {
	router := mux.NewRouter()
	sessionStore := sessions.NewCookieStore([]byte("something-very-secret"))
	configureRoutes(router, sessionStore)

	http.Handle("/", router)
	http.ListenAndServe(config.Settings.Port, context.ClearHandler(http.DefaultServeMux))
}
