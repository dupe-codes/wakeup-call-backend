package main

import (
	"net/http"
	"fmt"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/njdup/wakeup-call-backend/conf"

	"./controllers/userController"
	"./models/group"
	"github.com/njdup/wakeup-call-backend/models/user"
)

func testGroupStuff(res http.ResponseWriter, req *http.Request) {

	user, err := user.FindMatchingUser("njdup")
    newGroup := &group.Group {
        Name: "TestGroup11",
    }
    err = newGroup.Save()
    if err != nil {
        fmt.Fprintf(res, "Error occurred saving the group: %s", err.Error())
        return
    }

    err = newGroup.AddUser(user)
    if err != nil {
        fmt.Fprintf(res, "Error adding user to the group: %s", err.Error())
        return
    }

    // Grab group again to update its list of users
    newGroup, err = group.FindMatchingGroup(newGroup.Name)
    if err != nil {
        fmt.Fprintf(res, "Error finding group with matching name: %s", err.Error())
        return
    }

    users, err := newGroup.Users()
    if err != nil {
        fmt.Fprintf(res, "Error encountered getting group users: %s", err.Error())
        return
    }
    userNames := []string{}
    for _, user := range users {
        userNames = append(userNames, user.Username)
    }
    // users shouldn't be an empty slice now
    returnString := "Group %s successfully created with following users: " + strings.Join(userNames, ", ")
    fmt.Fprintf(res, returnString, newGroup.Name)
    return
}

// ConfigureRoutes sets all API routes
func configureRoutes(router *mux.Router, sessionStore *sessions.CookieStore) {
	userController.ConfigRoutes(router, sessionStore)
	router.HandleFunc("/test", testGroupStuff)
}

// Main launches the API server
func main() {
	router := mux.NewRouter()
	sessionStore := sessions.NewCookieStore([]byte("something-very-secret"))
	configureRoutes(router, sessionStore)

	http.Handle("/", router)
	http.ListenAndServe(config.Settings.Port, context.ClearHandler(http.DefaultServeMux))
}
