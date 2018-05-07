package acmogo_test

import (
	"testing"

	"github.com/crhntr/acmogo"
	"github.com/globalsign/mgo/bson"
)

func TestACL(t *testing.T) {
	user0 := User{Entity: acmogo.New()}
	user1 := User{Entity: acmogo.New()}
	// user2 := User{Entity: acmogo.New()}

	team0 := Team{Entity: acmogo.New()}
	team1 := Team{Entity: acmogo.New()}
	post0 := Post{Entity: acmogo.New()}

	idacentityZeroVal := acmogo.Reference{}
	if err := post0.SetCreator(idacentityZeroVal); err == nil {
		t.Fatal()
	}
	if err := post0.SetCreator(user0.Ref()); err != nil {
		t.Fatal()
	}
	if err := post0.SetCreator(user0.Ref()); err == nil {
		t.Fatal()
	}

	if !post0.DeletePermitted(user0.Ref()) {
		t.Fatal()
	}
	if !post0.UpdatePermitted(user0.Ref()) {
		t.Fatal()
	}
	if !post0.ReadPermitted(user0.Ref()) {
		t.Fatal()
	}

	if post0.ReadPermitted(user1.Ref()) {
		t.Fatal()
	}
	if post0.UpdatePermitted(user1.Ref()) {
		t.Fatal()
	}
	if post0.DeletePermitted(user1.Ref()) {
		t.Fatal()
	}

	if !post0.ReadPermitted(user1.Ref(), user0.Ref()) {
		t.Fatal()
	}

	if !post0.ReadPermitted(user0.Ref(), user1.Ref()) {
		t.Fatal()
	}

	post0.PermitDelete(team0.Ref(), team1.Ref())

	post0.ClearAccessControl(user1.Ref())
	if post0.DeletePermitted(user1.Ref()) {
		t.Fatal()
	}
	if post0.UpdatePermitted(user1.Ref()) {
		t.Fatal()
	}
	if post0.ReadPermitted(user1.Ref()) {
		t.Fatal()
	}

	post0.PermitRead(user1.Ref())
	if post0.DeletePermitted(user1.Ref()) {
		t.Fatal()
	}
	if post0.UpdatePermitted(user1.Ref()) {
		t.Fatal()
	}
	if !post0.ReadPermitted(user1.Ref()) {
		t.Fatal()
	}

	post0.PermitUpdate(user1.Ref())
	if post0.DeletePermitted(user1.Ref()) {
		t.Fatal()
	}
	if !post0.UpdatePermitted(user1.Ref()) {
		t.Fatal()
	}
	if !post0.ReadPermitted(user1.Ref()) {
		t.Fatal()
	}

	post0.PermitDelete(user1.Ref())
	if !post0.DeletePermitted(user1.Ref()) {
		t.Fatal()
	}
	if !post0.UpdatePermitted(user1.Ref()) {
		t.Fatal()
	}
	if !post0.ReadPermitted(user1.Ref()) {
		t.Fatal()
	}

	post0.PermitUpdate(user1.Ref())
	if post0.DeletePermitted(user1.Ref()) {
		t.Fatal()
	}
	if !post0.UpdatePermitted(user1.Ref()) {
		t.Fatal()
	}
	if !post0.ReadPermitted(user1.Ref()) {
		t.Fatal()
	}

	post0.PermitRead(user1.Ref())
	if post0.DeletePermitted(user1.Ref()) {
		t.Fatal()
	}
	if post0.UpdatePermitted(user1.Ref()) {
		t.Fatal()
	}
	if !post0.ReadPermitted(user1.Ref()) {
		t.Fatal()
	}

	post0.ClearAccessControl(user1.Ref())
	if post0.DeletePermitted(user1.Ref()) {
		t.Fatal()
	}
	if post0.UpdatePermitted(user1.Ref()) {
		t.Fatal()
	}
	if post0.ReadPermitted(user1.Ref()) {
		t.Fatal()
	}

	post0.Public = true
	user2 := User{Entity: acmogo.New()}
	if !post0.ReadPermitted(user2) {
		t.Error("read should be permitted")
	}
}

func TestMakeReferenceList(t *testing.T) {
	acmogo.MakeReferenceList("col", []bson.ObjectId{bson.NewObjectId(), bson.NewObjectId()}...)
}
