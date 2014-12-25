package user

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    "github.com/njdup/wakeup-call-backend/conf"
    "../../db/"
)

type User struct {
    Id  bson.ObjectId   `bson:"_id, omitempty" json: "-"`
    Username string     `bson:"userName" json:"userName"`
    Fullname string     `bson:"fullName" json:"fullName"`
    PasswordHash string `bson:"passwordHash" json:"-"`
    PasswordSalt string `bson:"passwordSalt" json:"-"`
    Inserted time.Time  `bson:"inserted" json:"-"`

}

var (
    collectionName = "users"
)

func TestInsert() {
    testEntry := &User{
        Username: "test",
        Fullname: "test_guy",
        PasswordHash: "hello",
        PasswordSalt: "whatsup?"
    }

    query := func(coll *mgo.Collection) error {
        count, err := coll.Find(bson.M{"userName": testEntry.Username}).Limit(1).Count()
        if err != nil {
            return err
        }
        if count > 0 {
            return Error("Entry already exists")
        }
        return coll.Insert(testEntry)
    }

    err := db.ExecWithColl(collectionName, query)
    if err != nil {
        panic(err)
    }
    return
}
