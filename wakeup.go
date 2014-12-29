package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/njdup/wakeup-call-backend/conf"

	"./controllers/groupController"
	"./controllers/userController"
	"./models/group"
	"github.com/njdup/wakeup-call-backend/models/user"
)

func testGroupStuff(res http.ResponseWriter, req *http.Request) {

	aUser, err := user.FindMatchingUser("njdup")
	newGroup := &group.Group{
		//Name: "TestGroup19",
		Name: "TheTrynas",
	}
	err = newGroup.Save()
	if err != nil {
		fmt.Fprintf(res, "Error occurred saving the group: %s", err.Error())
		return
	}
	newGroup, err = group.FindMatchingGroup(newGroup.Name)

	err = newGroup.AddUser(aUser)
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

	users, err := newGroup.GetUsers()
	if err != nil {
		fmt.Fprintf(res, "Error encountered getting group users: %s", err.Error())
		return
	}
	userNames := []string{}
	for _, aUser := range users {
		userNames = append(userNames, aUser.Username)
	}
	// users shouldn't be an empty slice now
	//returnString := "Group %s successfully created with following users: " + strings.Join(userNames, ", ")
	//fmt.Fprintf(res, returnString, newGroup.Name)
	aUser, err = user.FindMatchingUser("njdup")
	userGroups, err := group.GetGroupsForUser(aUser)
	groupNames := []string{}
	for _, userGroup := range userGroups {
		groupNames = append(groupNames, userGroup.Name)
	}
	returnString := "User %s is now in the following groups: " + strings.Join(groupNames, ", ")
	fmt.Fprintf(res, returnString, aUser.Username)
	return
}

// ConfigureRoutes sets all API routes
func configureRoutes(router *mux.Router, sessionStore *sessions.CookieStore) {
	userController.ConfigRoutes(router, sessionStore)
	groupController.ConfigRoutes(router, sessionStore)
	router.HandleFunc("/test", testGroupStuff)
}

// Main launches the API server
func main() {
	router := mux.NewRouter()
	sessionStore := sessions.NewCookieStore([]byte("something-very-secret")) //TODO: Fix secret key
	configureRoutes(router, sessionStore)

	http.Handle("/", router)
	http.ListenAndServe(config.Settings.Port, context.ClearHandler(http.DefaultServeMux))
}
