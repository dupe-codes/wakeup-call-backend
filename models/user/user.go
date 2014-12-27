package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/njdup/wakeup-call-backend/conf"
	"github.com/njdup/wakeup-call-backend/db"
	"github.com/njdup/wakeup-call-backend/utils/errors"
	"github.com/njdup/wakeup-call-backend/utils/security"
)

type User struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json: "-"`
	Username     string        `bson:"userName" json:"userName"`
	Fullname     string        `bson:"fullName" json:"fullName"`
	PasswordHash string        `bson:"passwordHash" json:"-"`
	PasswordSalt string        `bson:"passwordSalt" json:"-"`
	Inserted     time.Time     `bson:"inserted" json:"-"`
}

var (
	collectionName = "users"
)

// ToString returns a string representation of the receiving user
func (user *User) ToString() string {
	return fmt.Sprintf("User %s: %s", user.Username, user.Fullname)
}

// Save inserts the receiver User into the database.
// Returns an error if one is encountered, including
// validation errors such as a user with the set username
// already existing
func (user *User) Save() error {
	emptyFields := checkEmptyFields(user)
	if len(emptyFields) != 0 {
		invalid := strings.Join(emptyFields, " ")
		return &errorUtils.InvalidFieldsError{
			"The following fields cannot be empty: " + invalid,
			emptyFields,
		}
	}

	insertQuery := func(col *mgo.Collection) error {
		count, err := col.Find(bson.M{"userName": user.Username}).Limit(1).Count()
		if err != nil {
			return err
		}
		if count != 0 {
			return &errorUtils.InvalidFieldsError{
			    "A user with the given username already exists",
			    []string{"Username"},
			}
		}
		user.Inserted = time.Now() // Add insertion time stamp
		return col.Insert(user)
	}

	return db.ExecWithCol(collectionName, insertQuery)
}

// HashPassword hashes the given password and saves it in the user struct
func (user *User) HashPassword(password string) error {
    passwordSalt := security.GenerateSalt()
    hashedPass := security.RunSHA2(password + passwordSalt)

    user.PasswordHash = hashedPass
    user.PasswordSalt = passwordSalt
    return nil
}

// ConfirmPassword checks if the given password matches the saved password
func (user *User) ConfirmPassword(givenPass string) bool {
    return security.RunSHA2(givenPass + user.PasswordSalt) == user.PasswordHash
}

/*
 * User model utility functions
 */

// checkEmptyFields ensures all required fields of a user obj are set
func checkEmptyFields(user *User) []string {
	result := make([]string, 0)

	// TODO: Find a better way of doing this kind of checking
	// For some reason, request of form Username: "" was successful. Look into 
	// this.
	if user.Username == "" {
		result = append(result, "username")
	}
	return result
}

/*
 * Testing Code (Bad form, should just make actual tests, but this'll do as I
 * continue learning Go to begin with ^_^)
 */

func TestInsert() {
	testEntry := &User{
		Username:     "test",
		Fullname:     "test_guy",
		PasswordHash: "hello",
		PasswordSalt: config.Settings.DatabaseUrl,
	}

	query := func(coll *mgo.Collection) error {
		count, err := coll.Find(bson.M{"userName": testEntry.Username}).Limit(1).Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("Entry already exists")
		}
		return coll.Insert(testEntry)
	}

	err := db.ExecWithCol(collectionName, query)
	if err != nil {
		panic(err)
	}
	return
}
