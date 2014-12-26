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
    errors []error
}

func AllUsers(res http.ResponseWriter, req *http.Request) {
	return
}

func GetUser(res http.ResponseWriter, req *http.Request) {
	return
}

func CreateUser(res http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(res, "This is where we'll create a new user, yay!")
	decoder := json.NewDecoder(req.Body)
	var newUser user.User
	err := decoder.Decode(&newUser)
	fmt.Fprintf(res, newUser.ToString())
	return
}

func Login(res http.ResponseWriter, req *http.Request) {
	return
}

func Logout(res http.ResponseWriter, req *http.Request) {
	return
}
