package db

import (
	"gopkg.in/mgo.v2"

	"github.com/njdup/wakeup-call-backend/conf"
)

type queryFunc func(*mgo.Collection) error

var (
	mgoSession *mgo.Session
)

func getDbSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(config.Settings.DatabaseUrl)
		if err != nil {
			panic(err) // TODO: Figure out better error handling
		}
	}
	return mgoSession.Clone()
}

// CloseDbSession closes the global db connection
// TODO: Needed?
func closeDbSession() {
	if mgoSession != nil {
		mgoSession.Close()
	}
}

// ExecWithCol executes the given query function on the given database
// collection
// Returns any error encountered from executing the query function
func ExecWithCol(collection string, fn queryFunc) error {
	session := getDbSession()
	defer session.Close()
	col := session.DB(config.Settings.DatabaseName).C(collection)
	return fn(col)
}
