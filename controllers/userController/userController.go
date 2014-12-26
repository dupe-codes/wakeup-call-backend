package userController

import (
	"encoding/json"
	"fmt"
	"net/http"

	//"gopkg.in/mgo.v2/bson"

	//"github.com/njdup/wakeup-call-backend/conf"
	"github.com/njdup/wakeup-call-backend/models/users"
)

type UserAPIResponse struct {
	Status int
	Data   interface{}
	Error  error
}

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

	// Now attempt to save, create appropriate response
	resContent := &UserAPIResponse{}
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
		panic(err)
	}
	fmt.Fprintf(res, string(payload))
	return
}

func Login(res http.ResponseWriter, req *http.Request) {
	return
}

func Logout(res http.ResponseWriter, req *http.Request) {
	return
}
