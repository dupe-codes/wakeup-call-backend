package main

import (
    "net/http"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/njdup/wakeup-call-backend/conf"
)

func HomeHandler(res http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(res, "This is a test!")
}

// Main configures all API routes and runs the server on the port specified in
// the project settings
func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", HomeHandler)
    http.Handle("/", router)

    http.ListenAndServe(config.Settings.Port, nil)
}
