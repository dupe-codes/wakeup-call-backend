package groupController

import (
    "fmt"
    "net/http"
    "encoding/json"

    "github.com/gorilla/sessions"
    "github.com/gorilla/mux"

    "github.com/njdup/wakeup-call-backend/models/group"
    "github.com/njdup/wakeup-call-backend/models/user"
    "github.com/njdup/wakeup-call-backend/utils/responses"
    "github.com/njdup/wakeup-call-backend/utils/errors"
)

func CreateGroup(sessionStore *sessions.CookieStore) http.Handler {
    return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        session, _ := sessionStore.Get(req, "wakeup-session")
        username, authenticated := session.Values["user"]
        if !authenticated {
            errorMsg := &errorUtils.GeneralError{Message: "You must sign in to create a group"}
            APIResponses.SendErrorResponse(errorMsg, http.StatusBadRequest, res)
            return
            /*
            resContent := &APIResponses.Response{Status: http.StatusBadRequest, Error: errorMsg}
            response, err := json.MarshalIndent(resContent, "", "  ")
            if err != nil {
                http.Error(res, "Error preparing response", http.StatusInternalServerError)
                return
            }
            http.Error(res, string(response), http.StatusBadRequest)
            return
            */
        }

        // Create and save the new group
        req.ParseForm()
        newGroup := &group.Group{Name: req.PostFormValue("Name")}
        err := newGroup.Save()
        if err != nil {
            // TODO: Repeated code here, refactor into responses util function
            resContent := &APIResponses.Response{Status: http.StatusBadRequest, Error: err}
            response, err := json.MarshalIndent(resContent, "", "  ")
            if err != nil {
                http.Error(res, "Error preparing response", http.StatusInternalServerError)
                return
            }
            http.Error(res, string(response), http.StatusBadRequest)
            return
        }

        // Now auto add creator to the new group
        newGroup, err = group.FindMatchingGroup(newGroup.Name) // Need to update w/ db info
        usernameString := fmt.Sprintf("%s", username)
        creator, err := user.FindMatchingUser(usernameString)
        err = newGroup.AddUser(creator)

        // If we reach here, group has been created + user added successfully
        // TODO: Make SuccessResponse util function
        resContent := &APIResponses.Response{Status: http.StatusOK, Data: "Success"}
        response, err := json.MarshalIndent(resContent, "", "  ")
        if err != nil {
            http.Error(res, "Error preparing response", http.StatusInternalServerError)
            return
        }
        fmt.Fprintf(res, string(response))
        return
    })
}

func ConfigRoutes(router *mux.Router, sessionStore *sessions.CookieStore) {
    router.Handle("/groups", CreateGroup(sessionStore)).Methods("POST")
}
