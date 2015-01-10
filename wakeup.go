package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	"github.com/njdup/wakeup-call-backend/conf"
	"github.com/njdup/wakeup-call-backend/controllers/groupController"
	"github.com/njdup/wakeup-call-backend/controllers/userController"
)

var (
	keyLength = 12
)

// ConfigureRoutes sets all API routes
func configureRoutes(router *mux.Router, sessionStore *sessions.CookieStore) {
	userController.ConfigRoutes(router, sessionStore)
	groupController.ConfigRoutes(router, sessionStore)
}

// Main launches the API server
func main() {
	router := mux.NewRouter()
	sessionStore := sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(keyLength)))
	configureRoutes(router, sessionStore)

	http.Handle("/", router)
	fmt.Println("Listening on port " + config.Settings.Port)
	http.ListenAndServe(config.Settings.Port, context.ClearHandler(http.DefaultServeMux))
}
