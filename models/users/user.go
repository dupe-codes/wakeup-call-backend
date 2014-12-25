package user

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../../db/"
	"github.com/njdup/wakeup-call-backend/conf"
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
