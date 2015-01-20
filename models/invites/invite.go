package invites

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/njdup/wakeup-call-backend/db"
	"github.com/njdup/wakeup-call-backend/utils/errors"
)

// TODO: Consider adding more robust features to invite system, such as invite codes
type Invite struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json"-"`
	Name        string        `bson:"name" json:"name"`
	Groupname   string        `bson:"groupName" json:"groupName"`
	Phonenumber string        `bson:"phoneNumber" json:"phoneNumber"`
	Created     time.Time     `bson:"created" json:"-"`
}

var (
	CollectionName = "invites"
)

func (invite *Invite) Save() error {
	// Add validation checks here
	emptyFields := checkEmptyFields(invite)
	if len(emptyFields) != 0 {
		invalid := strings.Join(emptyFields, " ")
		return &errorUtils.InvalidFieldsError{
			"The following fields cannot be empty: " + invalid,
			emptyFields,
		}
	}

	insertQuery := func(col *mgo.Collection) error {
		count, err := col.Find(bson.M{
			"groupName":   invite.Groupname,
			"phoneNumber": invite.Phonenumber,
		}).Limit(1).Count()
		if err != nil {
			return err
		}
		if count != 0 {
			return &errorUtils.InvalidFieldsError{
				"An invite for this group has already been made for the given number",
				[]string{"Phonenumber"},
			}
		}

		invite.Created = time.Now()
		return col.Insert(invite)
	}

	return db.ExecWithCol(CollectionName, insertQuery)
}

/*
 * Invite model utility functions
 */

// checkEmptyFields ensures all required fields of an invite obj are set
func checkEmptyFields(invite *Invite) []string {
	result := make([]string, 0)

	if invite.Name == "" {
		result = append(result, "Name")
	}
	if invite.Groupname == "" {
		result = append(result, "Groupname")
	}
	if invite.Phonenumber == "" {
		result = append(result, "Phonenumber")
	}
	return result
}
