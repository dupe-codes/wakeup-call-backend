package userController

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	//"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/sessions"

	//"github.com/njdup/wakeup-call-backend/conf"
	"github.com/njdup/wakeup-call-backend/models/user"
	"github.com/njdup/wakeup-call-backend/utils/responses"
)

var sessionStore = sessions.NewCookieStore([]byte("something-very-secret"))

func AllUsers(res http.ResponseWriter, req *http.Request) {
	// Simple test to see if user is authenticated
	session, _ := sessionStore.Get(req, "wakeup-session")
	username, authenticated := session.Values["user"]
	if !authenticated {
		http.Error(res, "You must sign in to gain access here", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(res, "Success! Welcome %s", username)
	return
}

func GetUser(res http.ResponseWriter, req *http.Request) {
	return
}

func CreateUser(res http.ResponseWriter, req *http.Request) {
	// Prepare new user from form data
	req.ParseForm()
	newUser := &user.User{
		Username: req.PostFormValue("Username"),
		Fullname: req.PostFormValue("Fullname"),
	}

	err := newUser.HashPassword(req.PostFormValue("Password"))
	if err != nil {
		fmt.Fprintf(res, "Password is invalid") //TODO: Handle with appropriate http resp
		return
	}

	// Now attempt to save, create appropriate response
	resContent := &APIResponses.Response{}
	err = newUser.Save()
	if err != nil {
		resContent.Status = 400
		resContent.Error = err
	} else {
		resContent.Status = 200
		resContent.Data = "success"
	}

	payload, err := json.MarshalIndent(resContent, "", "  ")
	if err != nil {
		fmt.Fprintf(res, `{"Status": 500, "Error": "Unable to prepare server response"}`)
		return
	}
	fmt.Fprintf(res, string(payload))
	return
}

func Login(res http.ResponseWriter, req *http.Request) {
	// First check if a session is already active
	session, _ := sessionStore.Get(req, "wakeup-session")
	if _, ok := session.Values["user"]; ok {
		http.Error(res, "User is already signed in", http.StatusBadRequest)
		return
	}

	// Otherwise, try authenticating
	req.ParseForm()
	matchedUser, err := user.FindMatchingUser(req.PostFormValue("Username"))
	if err != nil {
		fmt.Fprintf(res, err.Error())
		return
	}

	if matchedUser.ConfirmPassword(req.PostFormValue("Password")) {
		session.Values["user"] = matchedUser.Username
		session.Save(req, res)

		resContent := &APIResponses.Response{Status: 200, Data: "Successfully signed in"}
		payload, _ := json.MarshalIndent(resContent, "", "  ")
		fmt.Fprintf(res, string(payload))
	} else {
		err = errors.New("Given password is incorrect")
		resContent := &APIResponses.Response{Status: 400, Error: err}
		payload, _ := json.MarshalIndent(resContent, "", " ")
		fmt.Fprintf(res, string(payload))
	}
	return
}

func Logout(res http.ResponseWriter, req *http.Request) {
	session, _ := sessionStore.Get(req, "wakeup-session")
	if _, ok := session.Values["user"]; !ok {
		http.Error(res, "No user present to sign out", http.StatusBadRequest)
		return
	}
	delete(session.Values, "user")
	session.Save(req, res)

	resContent := &APIResponses.Response{Status: 200, Data: "Successfully logged out"}
	payload, _ := json.MarshalIndent(resContent, "", "  ")
	fmt.Fprintf(res, string(payload))
	return
}
