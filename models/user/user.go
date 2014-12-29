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
	"github.com/njdup/wakeup-call-backend/models/group"
)

type User struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json: "-"`
	Username     string        `bson:"userName" json:"userName"`
	Fullname     string        `bson:"fullName" json:"fullName"`
	PasswordHash string        `bson:"passwordHash" json:"-"`
	PasswordSalt string        `bson:"passwordSalt" json:"-"`
	Inserted     time.Time     `bson:"inserted" json:"-"`

	Groups []bson.ObjectID `bson:"groups" json:"groups"`
}

var (
	CollectionName = "users"
)

// ToString returns a string representation of the receiving user
func (user *User) ToString() string {
	return fmt.Sprintf("User %s: %s with (pass: %s, salt: %s)", user.Username, user.Fullname, user.PasswordHash, user.PasswordSalt)
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

	return db.ExecWithCol(CollectionName, insertQuery)
}

// HashPassword hashes the given password and saves it in the user struct
func (user *User) HashPassword(password string) error {
	if !passwordValid(password) {
		return &errorUtils.InvalidFieldsError{
			"The given password is not acceptable",
			[]string{"Password"},
		}
	}
	passwordSalt := security.GenerateSalt()
	hashedPass := security.RunSHA2(password + passwordSalt)

	user.PasswordHash = hashedPass
	user.PasswordSalt = passwordSalt
	return nil
}

// ConfirmPassword checks if the given password matches the saved password
func (user *User) ConfirmPassword(givenPass string) bool {
	return security.RunSHA2(givenPass+user.PasswordSalt) == user.PasswordHash
}

// FindMatchingUser searches for a saved user with the given username
// Returns a pointer to the matching User struct
// TODO: Figure out how to detect when no matching user found
func FindMatchingUser(username string) (*User, error) {
	result := User{}
	searchQuery := func(col *mgo.Collection) error {
		return col.Find(bson.M{"userName": username}).One(&result)
	}
	err := db.ExecWithCol(CollectionName, searchQuery)
	if err != nil {
		return nil, err
	}
	return &result, nil
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

// passwordValid checks if the password conforms to the password policy
// TODO: Have this return a PasswordValidationError rather than a bool
// Maybe have validators return the error with appropriate message
func passwordValid(password string) bool {
	for _, validator := range security.PasswordPolicy.Validations {
		if !validator(password) {
			return false
		}
	}
	return true
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

	err := db.ExecWithCol(CollectionName, query)
	if err != nil {
		panic(err)
	}
	return
}
