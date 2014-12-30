// TODO: Clean up this file with new tricks like your new response utils

package userController

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/njdup/wakeup-call-backend/models/user"
	"github.com/njdup/wakeup-call-backend/utils/errors"
	"github.com/njdup/wakeup-call-backend/utils/responses"
)

func AllUsers(sessionStore *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		return
	})
}

func GetUser(res http.ResponseWriter, req *http.Request) {
	return
}

func CreateUser(sessionStore *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Prepare new user from form data
		req.ParseForm()
		newUser := &user.User{
			Username:  req.FormValue("Username"),
			Firstname: req.FormValue("Firstname"),
			Lastname:  req.FormValue("Lastname"),
		}

		// TODO: Add parsePhonenumber to user model
		phonenumber, err := parsePhonenumber(req.FormValue("Phonenumber"))
		if err != nil {
			errorMsg := &errorUtils.InvalidFieldsError{
				Message: "Given phone number is invalid",
				Fields:  []string{"Phonenumber"},
			}
			APIResponses.SendErrorResponse(errorMsg, http.StatusBadRequest, res)
			return
		}
		newUser.Phonenumber = phonenumber

		err = newUser.HashPassword(req.FormValue("Password"))
		if err != nil {
			errorMsg := &errorUtils.InvalidFieldsError{
				Message: "The given password is invalid",
				Fields:  []string{"Password"},
			}
			APIResponses.SendErrorResponse(errorMsg, http.StatusBadRequest, res)
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
	})
}

func Login(sessionStore *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// First check if a session is already active
		session, _ := sessionStore.Get(req, "wakeup-session")
		if _, ok := session.Values["user"]; ok {
			http.Error(res, "User is already signed in", http.StatusBadRequest)
			return
		}

		// Otherwise, try authenticating
		req.ParseForm()
		matchedUser, err := user.FindMatchingUser(req.FormValue("Username"))
		if err != nil {
			fmt.Fprintf(res, err.Error())
			return
		}

		if matchedUser.ConfirmPassword(req.FormValue("Password")) {
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
	})
}

func Logout(sessionStore *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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
	})
}

// ConfigRoutes initializes all application routes specific to users
func ConfigRoutes(router *mux.Router, sessionStore *sessions.CookieStore) {
	router.Handle("/users", CreateUser(sessionStore)).Methods("POST")
	router.Handle("/users/login", Login(sessionStore)).Methods("POST")
	router.Handle("/users/logout", Logout(sessionStore)).Methods("POST")
}

// TODO: Write this to format phone numbers given in form
func parsePhonenumber(inputNumber string) (string, error) {
	return inputNumber, nil
}
