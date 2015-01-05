package groupController

import (
	"fmt"
	"net/http"
	//"encoding/json"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/njdup/wakeup-call-backend/models/group"
	"github.com/njdup/wakeup-call-backend/models/user"
	"github.com/njdup/wakeup-call-backend/utils/errors"
	"github.com/njdup/wakeup-call-backend/utils/responses"
)

func CreateGroup(sessionStore *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		session, _ := sessionStore.Get(req, "wakeup-session")
		username, authenticated := session.Values["user"]
		if !authenticated {
			errorMsg := &errorUtils.GeneralError{Message: "You must sign in to create a group"}
			APIResponses.SendErrorResponse(errorMsg, http.StatusBadRequest, res)
			return
		}

		// Create and save the new group
		req.ParseForm()
		newGroup := &group.Group{Name: req.FormValue("Name")}
		err := newGroup.Save()
		if err != nil {
			APIResponses.SendErrorResponse(err, http.StatusBadRequest, res)
			return
		}

		// After group created, provision its phone number
		err = newGroup.ProvisionPhoneNumber()
		if err != nil {
            errorMsg := &errorUtils.GeneralError{"Error provisioning group phone number"}
            APIResponses.SendErrorResponse(errorMsg, http.StatusInternalServerError, res)
            return
		}

		// Now auto add creator to the new group
		// TODO: Handle errors
		newGroup, err = group.FindMatchingGroup(newGroup.Name) // Need to update w/ db info
		usernameString := fmt.Sprintf("%s", username)
		creator, err := user.FindMatchingUser(usernameString)
		err = newGroup.AddUser(creator)

		// If we reach here, group has been created + user added successfully
		APIResponses.SendSuccessResponse("Success", res)
		return
	})
}

func GetGroupUsers(sessionStore *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		groupName := vars["groupName"]
		group, err := group.FindMatchingGroup(groupName)
		if err != nil {
			errorMsg := &errorUtils.GeneralError{"No group matching the given group name was found"}
			APIResponses.SendErrorResponse(errorMsg, http.StatusBadRequest, res)
			return
		}

		users, err := group.GetUsers()
		if err != nil {
			errorMsg := &errorUtils.GeneralError{"Error occurred fetching the group's users"}
			APIResponses.SendErrorResponse(errorMsg, http.StatusInternalServerError, res)
			return
		}

		APIResponses.SendSuccessResponse(users, res)
		return
	})
}

func GetGroup(sessionStore *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		groupName := vars["groupName"]
		group, err := group.FindMatchingGroup(groupName)
		if err != nil {
			errorMsg := &errorUtils.GeneralError{"No group matching the given group name was found"}
			APIResponses.SendErrorResponse(errorMsg, http.StatusBadRequest, res)
			return
		}

		APIResponses.SendSuccessResponse(group, res)
		return
	})
}

func ConfigRoutes(router *mux.Router, sessionStore *sessions.CookieStore) {
	router.Handle("/groups", CreateGroup(sessionStore)).Methods("POST")
	router.Handle("/groups/{groupName}/users", GetGroupUsers(sessionStore)).Methods("GET")
	router.Handle("/groups/{groupName}", GetGroup(sessionStore)).Methods("GET")
}
