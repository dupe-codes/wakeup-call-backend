package test

import (
    "net/http"
    "fmt"
)

func TestHandler(res http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(res, "This is a test of putting handlers in separate packages...")
}
