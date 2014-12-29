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
		newGroup := &group.Group{Name: req.PostFormValue("Name")}
		err := newGroup.Save()
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

func ConfigRoutes(router *mux.Router, sessionStore *sessions.CookieStore) {
	router.Handle("/groups", CreateGroup(sessionStore)).Methods("POST")
}
