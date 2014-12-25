package user

import (
	"errors"
	"time"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/njdup/wakeup-call-backend/conf"
	"github.com/njdup/wakeup-call-backend/db"
)

type User struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json: "-"`
	Username     string        `bson:"userName" json:"userName"`
	Fullname     string        `bson:"fullName" json:"fullName"`
	PasswordHash string        `bson:"passwordHash" json:"-"`
	PasswordSalt string        `bson:"passwordSalt" json:"-"`
	Inserted     time.Time     `bson:"inserted" json:"-"`
}

type InvalidFieldsError struct {
    msg string
    Fields []string
}

func (err *InvalidFieldsError) Error() string { return err.msg }

var (
	collectionName = "users"
)

// Save inserts the receiver User into the database.
// Returns an error if one is encountered, including
// validation errors such as a user with the set username
// already existing
func (user *User) Save() error {
    emptyFields := checkEmptyFields(user)
    if len(emptyFields) != 0 {
        invalid := strings.Join(emptyFields, " ")
        return errors.New("The following fields cannot be empty: " + invalid)
    }

    insertQuery := func(col *mgo.Collection) error {
        count, err := col.Find(bson.M{"userName": user.Username}).Limit(1).Count()
        if err != nil {
            return err
        }
        if count != 0 {
            return &InvalidFieldsError{"A user with the given username already exists", []string{"Username"}}
        }
        return col.Insert(user)
    }

    // TODO: Add insert time stamp here

    return db.ExecWithCol(collectionName, insertQuery)
}


/*
 * User model utility functions
 */

// checkEmptyFields ensures all required fields of a user obj are set
func checkEmptyFields(user *User) []string {
    result := make([]string, 0)

    // TODO: Find a better way of doing this kind of checking
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
