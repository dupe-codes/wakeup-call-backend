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
    status int
    data interface{}
    err error
}

func AllUsers(res http.ResponseWriter, req *http.Request) {
	return
}

func GetUser(res http.ResponseWriter, req *http.Request) {
	return
}

func CreateUser(res http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(res, "This is where we'll create a new user, yay!")
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
        resContent.status = 400
        resContent.err = err
	} else {
        resContent.status = 200
        resContent.data = []byte{'h', 'i'}
	}
    payload, err := json.MarshalIndent(resContent, "", " ")
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
