package userController

import (
	"encoding/json"
	"fmt"
	"net/http"

	//"gopkg.in/mgo.v2/bson"

	//"github.com/njdup/wakeup-call-backend/conf"
	"github.com/njdup/wakeup-call-backend/models/user"
	"github.com/njdup/wakeup-call-backend/utils/responses"
)

func AllUsers(res http.ResponseWriter, req *http.Request) {
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
	// TODO: Take given password here and generate hash + salt.
	// Make helpers to do this
	newUser.HashPassword(req.PostFormValue("Password"))
	fmt.Printf(newUser.ToString())

	// Now attempt to save, create appropriate response
	resContent := &APIResponses.Response{}
	err := newUser.Save()
	if err != nil {
		resContent.Status = 400
		resContent.Error = err
	} else {
		resContent.Status = 200
		resContent.Data = "success"
	}

	payload, err := json.MarshalIndent(resContent, "", "  ")
	if err != nil {
		fmt.Fprintf(res, `{"Status": 400, "Error": "Unknown"}`)
		return
	}
	fmt.Fprintf(res, string(payload))
	return
}

func Login(res http.ResponseWriter, req *http.Request) {
	// Do authentication stuff to create new session here
	// TODO: Use helper function to check if hashed given pass +
	// salt == stored hash + salt
	return
}

func Logout(res http.ResponseWriter, req *http.Request) {
	// Do unauth stuff to destroy session here
	return
}
