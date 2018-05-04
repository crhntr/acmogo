package entity_test

import (
	"github.com/crhntr/litsphere/internal/entity"
	"github.com/globalsign/mgo/bson"
)

const (
	UserCol = "user"
	TeamCol = "team"
	PostCol = "post"
)

type (
	User struct {
		entity.Entity `bson:",inline"`
		Teams         []bson.ObjectId `bson:"teams"`
	}

	Team struct {
		entity.Entity `bson:",inline"`
	}

	Post struct {
		entity.Entity `bson:",inline"`
		N             int `bson:"n"`
	}
)

func (this User) Ref() entity.Reference {
	return entity.Reference{UserCol, this.ID}
}

func (this Team) Ref() entity.Reference {
	return entity.Reference{TeamCol, this.ID}
}

func (this Post) Ref() entity.Reference {
	return entity.Reference{PostCol, this.ID}
}
