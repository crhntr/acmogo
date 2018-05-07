package acmogo_test

import (
	"github.com/crhntr/acmogo"
	"github.com/globalsign/mgo/bson"
)

const (
	UserCol = "user"
	TeamCol = "team"
	PostCol = "post"
)

type (
	User struct {
		acmogo.Entity `bson:",inline"`
		Teams         []bson.ObjectId `bson:"teams"`
	}

	Team struct {
		acmogo.Entity `bson:",inline"`
	}

	Post struct {
		acmogo.Entity `bson:",inline"`
		N             int `bson:"n"`
	}
)

func (this User) Ref() acmogo.Reference {
	return acmogo.Reference{UserCol, this.ID}
}

func (this Team) Ref() acmogo.Reference {
	return acmogo.Reference{TeamCol, this.ID}
}

func (this Post) Ref() acmogo.Reference {
	return acmogo.Reference{PostCol, this.ID}
}
