package userController

import (
	//"encoding/json"
	"fmt"
	"net/http"

	//"gopkg.in/mgo.v2/bson"

	//"github.com/njdup/wakeup-call-backend/conf"
	//"github.com/njdup/wakeup-call-backend/models/users"
)

func AllUsers(res http.ResponseWriter, req *http.Request) {
	return
}

func GetUser(res http.ResponseWriter, req *http.Request) {
	return
}

func CreateUser(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "This is where we'll create a new user, yay!")
	return
}

func Login(res http.ResponseWriter, req *http.Request) {
	return
}

func Logout(res http.ResponseWriter, req *http.Request) {
	return
}
