package groupController

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/njdup/wakeup-call-backend/models/group"
	"github.com/njdup/wakeup-call-backend/models/invites"
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

		req.ParseForm()
		newGroup := &group.Group{Name: req.FormValue("Name")}
		err := newGroup.ProvisionPhoneNumber()
		if err != nil {
			APIResponses.SendErrorResponse(err, http.StatusInternalServerError, res)
			return
		}

		err = newGroup.Save()
		if err != nil {
			APIResponses.SendErrorResponse(err, http.StatusBadRequest, res)
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

func GetGroupInfo(sessionStore *sessions.CookieStore) http.Handler {
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

func GetGroup(sessionStore *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		queryValues := req.URL.Query()
		if len(queryValues["phoneNumber"]) == 0 {
			errorMsg := &errorUtils.GeneralError{Message: "A phone number must be included to query for a group"}
			APIResponses.SendErrorResponse(errorMsg, http.StatusBadRequest, res)
			return
		}

		phoneNumber := queryValues["phoneNumber"][0]
		fmt.Println("Received query for group with phone number: " + phoneNumber)
		group, err := group.FindGroupWithNumber(phoneNumber)
		// TODO: Improve this error handling
		if err != nil {
			errorMsg := &errorUtils.GeneralError{Message: "Unable to find group with the given phone number"}
			APIResponses.SendErrorResponse(errorMsg, http.StatusInternalServerError, res)
			return
		}
		APIResponses.SendSuccessResponse(group, res)
		return
	})
}

func CreateGroupInvite(sessionStore *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// First ensure that a user is logged in
		session, _ := sessionStore.Get(req, "wakeup-session")
		if _, ok := session.Values["user"]; !ok {
			errorMsg := &errorUtils.GeneralError{Message: "You must be signed in to invite a user to join a group"}
			APIResponses.SendErrorResponse(errorMsg, http.StatusBadRequest, res)
			return
		}

		req.ParseForm()
		vars := mux.Vars(req)
		groupName := vars["groupName"]

		// TODO: Add more robust validations, such as checking given phone number
		invite := &invites.Invite{
			Name:        req.FormValue("Name"),
			Groupname:   groupName,
			Phonenumber: req.FormValue("Phonenumber"),
		}

		// TODO: Add more robust error handling
		err := invite.Save()
		if err != nil {
			APIResponses.SendErrorResponse(err, http.StatusBadRequest, res)
			return
		}

		APIResponses.SendSuccessResponse("Invite successfully created", res)
		return
	})
}

func ConfigRoutes(router *mux.Router, sessionStore *sessions.CookieStore) {
	router.Handle("/groups", CreateGroup(sessionStore)).Methods("POST")
	router.Handle("/groups", GetGroup(sessionStore)).Methods("GET")
	router.Handle("/groups/{groupName}/users", GetGroupUsers(sessionStore)).Methods("GET")
	router.Handle("/groups/{groupName}", GetGroupInfo(sessionStore)).Methods("GET")
	router.Handle("/groups/{groupName}/invite", CreateGroupInvite(sessionStore)).Methods("POST")
}
