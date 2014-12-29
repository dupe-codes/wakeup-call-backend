package groupController

import (
    "fmt"
    "net/http"
    "encoding/json"

    "github.com/gorilla/sessions"
    "github.com/gorilla/mux"

    //"github.com/njdup/wakeup-call-backend/models/group"
    "github.com/njdup/wakeup-call-backend/utils/responses"
    "github.com/njdup/wakeup-call-backend/utils/errors"
)

func CreateGroup(sessionStore *sessions.CookieStore) http.Handler {
    return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        session, _ := sessionStore.Get(req, "wakeup-session")
        username, authenticated := session.Values["user"]
        if !authenticated {
            err := errorUtils.GeneralError{Message: "You must sign in to create a group"}
            resContent := &APIResponses.Response{Status: http.StatusBadRequest, Error: err}
            response, err := json.MarshalIndent(resContent, "", " ")
            if err != nil {
                http.Error(res, "Error preparing response", http.StatusInternalServerError)
                return
            }
            http.Error(res, string(response), http.StatusBadRequest)
            return
        }

        fmt.Fprintf(res, "%s is signed in", username)
        return
    })
}

func ConfigRoutes(router *mux.Router, sessionStore *sessions.CookieStore) {
    router.Handle("/groups", CreateGroup(sessionStore)).Methods("POST")
}
