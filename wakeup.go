package main

import (
    "net/http"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/njdup/wakeup-call-backend/conf"
    "./controllers/test"
)

func HomeHandler(res http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(res, "This is a test!")
}

// ConfigureRoutes sets all API routes
func configureRoutes(router *mux.Router) {
    router.HandleFunc("/test/", test.TestHandler)
}

// Main launches the API server
func main() {
    router := mux.NewRouter()
    configureRoutes(router)

    http.Handle("/", router)
    http.ListenAndServe(config.Settings.Port, nil)
}
