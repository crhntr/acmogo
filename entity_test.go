package acmogo_test

import (
	"flag"
	"os"
	"testing"

	"github.com/crhntr/acmogo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	dbName             = "mongox_entity_test"
	databaseSession, _ = mgo.DialWithInfo(&mgo.DialInfo{
		Database: dbName,
		Addrs:    []string{":27017"},
	})
)

type mp map[string]interface{}

var db *mgo.Database

func TestMain(m *testing.M) {
	flag.Parse()
	defer databaseSession.Close()
	db = databaseSession.DB("")
	os.Exit(m.Run())
}

func TestEntity(t *testing.T) {
	for name, test := range map[string]func(t *testing.T){
		"ReadPermitted": testEntityReadPermitted,
	} {
		t.Run(name, test)
		databaseSession.DB("").DropDatabase()
	}
}

func testEntityReadPermitted(t *testing.T) {

}

func TestEntity_Validate(t *testing.T) {
	identityZeroVal := acmogo.Reference{}
	if identityZeroVal.Validate() == nil {
		t.Fatal()
	}

	user0 := User{Entity: acmogo.New()}
	user1 := User{Entity: acmogo.New()}
	user2 := User{Entity: acmogo.New()}

	team0 := Team{Entity: acmogo.New()}
	team1 := Team{Entity: acmogo.New()}

	for _, idfr := range []acmogo.Referencer{user0, user1, user2, team0, team1} {
		if err := idfr.Ref().Validate(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestInsertEntity(t *testing.T) {
	user0 := User{Entity: acmogo.New()}
	user1 := User{Entity: acmogo.New()}
	team0 := Team{Entity: acmogo.New()}
	post0 := Post{Entity: acmogo.New()}

	if _, err := acmogo.InsertList(db, user0, user1, team0, post0); err != nil {
		t.Fatal(err)
	}
	if _, err := acmogo.InsertList(db, user0); err == nil {
		t.Fatal(err)
	}
}

func TestRefreshEntity(t *testing.T) {
	post0 := Post{Entity: acmogo.New()}
	acmogo.InsertList(db, post0)
	post0N := 101
	if err := db.C(PostCol).UpdateId(post0.ID, mp{"$set": mp{"n": post0N}}); err != nil {
		t.Fatal(err)
	}
	// t.Log(post0)
	acmogo.RefreshEntity(db, &post0)
	// t.Log(post0)
	if post0.N != post0N {
		t.Fail()
	}
}

func TestUpdateEntity(t *testing.T) {
	post0 := Post{Entity: acmogo.New()}
	acmogo.InsertList(db, post0)
	post0N := 101
	// t.Log(post0)
	acmogo.UpdateEntity(db, post0, bson.M{"$set": bson.M{"n": post0N}})
	acmogo.RefreshEntity(db, &post0)
	// t.Log(post0)
	if post0.N != post0N {
		t.Fail()
	}
}

func TestPersistClearAccessControl(t *testing.T) {
	post0 := Post{Entity: acmogo.New()}
	user0 := User{Entity: acmogo.New()}

	post0.PermitRead(user0)
	acmogo.InsertList(db, post0, user0)

	if err := acmogo.PersistClearAccessControl(db, post0, user0); err != nil {
		t.Fail()
	}
	if acmogo.DeletePermitted(db, post0, user0) {
		t.Fail()
	}
	if acmogo.UpdatePermitted(db, post0, user0) {
		t.Fail()
	}
	if acmogo.ReadPermitted(db, post0, user0) {
		t.Fail()
	}
}

func TestPersistPermitRead(t *testing.T) {
	post0 := Post{Entity: acmogo.New()}
	post1 := Post{Entity: acmogo.New()}
	user0 := User{Entity: acmogo.New()}
	user1 := User{Entity: acmogo.New()}

	acmogo.InsertList(db, post0, user0, user1)

	if acmogo.ReadPermitted(db, post0, user0, user1) {
		t.Error("read should not be permitted")
	}

	acmogo.PersistPermitRead(db, post0, user0)

	if !acmogo.ReadPermitted(db, post0, user0) {
		t.Error("read should be permitted")
	}
	if acmogo.UpdatePermitted(db, post0, user0) {
		t.Error("update should not be permitted")
	}
	if acmogo.DeletePermitted(db, post0, user0) {
		t.Error("delete should not be permitted")
	}
	if acmogo.ReadPermitted(db, post0, user1) {
		t.Error("read should not be permitted for unknown reader")
	}
	if acmogo.ReadPermitted(db, post1, user1) {
		t.Error("post1 should not be found")
	}
}

func TestPersistPermitUpdate(t *testing.T) {
	post0 := Post{Entity: acmogo.New()}
	post1 := Post{Entity: acmogo.New()}
	user0 := User{Entity: acmogo.New()}
	user1 := User{Entity: acmogo.New()}

	acmogo.InsertList(db, post0, user0, user1)

	if acmogo.UpdatePermitted(db, post0, user0, user1) {
		t.Error("read should not be permitted")
	}

	acmogo.PersistPermitUpdate(db, post0, user0)

	if !acmogo.ReadPermitted(db, post0, user0) {
		t.Error("read should be permitted")
	}
	if !acmogo.UpdatePermitted(db, post0, user0) {
		t.Error("update should be permitted")
	}
	if acmogo.DeletePermitted(db, post0, user0) {
		t.Error("delete should not be permitted")
	}

	if acmogo.UpdatePermitted(db, post0, user1) {
		t.Error("update should not be permitted for unknown updater")
	}
	if acmogo.UpdatePermitted(db, post1, user1) {
		t.Error("post1 should not be found")
	}
}

func TestPersistPermitDelete(t *testing.T) {
	post0 := Post{Entity: acmogo.New()}
	post1 := Post{Entity: acmogo.New()}
	user0 := User{Entity: acmogo.New()}
	user1 := User{Entity: acmogo.New()}

	acmogo.InsertList(db, post0, user0, user1)

	if acmogo.DeletePermitted(db, post0, user0, user1) {
		t.Error("delete should not be permitted")
	}

	acmogo.PersistPermitDelete(db, post0, user0)

	if !acmogo.ReadPermitted(db, post0, user0) {
		t.Error("read should be permitted")
	}
	if !acmogo.UpdatePermitted(db, post0, user0) {
		t.Error("update should be permitted")
	}
	if !acmogo.DeletePermitted(db, post0, user0) {
		t.Error("delete should be permitted")
	}

	if acmogo.DeletePermitted(db, post0, user1) {
		t.Error("delete should not be permitted for unknown deleter")
	}
	if acmogo.DeletePermitted(db, post1, user1) {
		t.Fail()
	}
}

func TestPersistPermitDeleteDowngradeToPermitUpdate(t *testing.T) {
	post0 := Post{Entity: acmogo.New()}
	user0 := User{Entity: acmogo.New()}

	acmogo.InsertList(db, post0, user0)

	acmogo.PersistPermitDelete(db, post0, user0)
	acmogo.PersistPermitUpdate(db, post0, user0)

	if !acmogo.ReadPermitted(db, post0, user0) {
		t.Fail()
	}
	if !acmogo.UpdatePermitted(db, post0, user0) {
		t.Fail()
	}
	if acmogo.DeletePermitted(db, post0, user0) {
		t.Fail()
	}
}

func TestPersistPermitDeleteDowngradeToPermitRead(t *testing.T) {
	post0 := Post{Entity: acmogo.New()}
	user0 := User{Entity: acmogo.New()}

	acmogo.InsertList(db, post0, user0)

	acmogo.PersistPermitDelete(db, post0, user0)
	acmogo.PersistPermitRead(db, post0, user0)

	if !acmogo.ReadPermitted(db, post0, user0) {
		t.Fail()
	}
	if acmogo.UpdatePermitted(db, post0, user0) {
		t.Fail()
	}
	if acmogo.DeletePermitted(db, post0, user0) {
		t.Fail()
	}
}

func TestPersistPermitUpdateDowngradeToPermitRead(t *testing.T) {
	post0 := Post{Entity: acmogo.New()}
	user0 := User{Entity: acmogo.New()}

	acmogo.InsertList(db, post0, user0)

	acmogo.PersistPermitUpdate(db, post0, user0)
	acmogo.PersistPermitRead(db, post0, user0)

	if !acmogo.ReadPermitted(db, post0, user0) {
		t.Fail()
	}
	if acmogo.UpdatePermitted(db, post0, user0) {
		t.Fail()
	}
	if acmogo.DeletePermitted(db, post0, user0) {
		t.Fail()
	}
}

func TestPersistPublic(t *testing.T) {
	post0 := Post{Entity: acmogo.New()}
	user0 := User{Entity: acmogo.New()}
	user1 := User{Entity: acmogo.New()}

	acmogo.InsertList(db, post0, user0, user1)

	if acmogo.ReadPermitted(db, post0, user0, user1) {
		t.Error("read should not be permitted")
	}

	acmogo.PersistPublic(db, post0)

	if acmogo.ReadPermitted(db, post0, user0, user1) {
		t.Error("read should be permitted")
	}
}

func TestPersistPrivate(t *testing.T) {
	post0 := Post{Entity: acmogo.New()}
	user0 := User{Entity: acmogo.New()}
	user1 := User{Entity: acmogo.New()}

	acmogo.InsertList(db, post0, user0, user1)

	if acmogo.ReadPermitted(db, post0, user0, user1) {
		t.Error("read should be permitted")
	}

	acmogo.PersistPrivate(db, post0)

	if acmogo.ReadPermitted(db, post0, user0, user1) {
		t.Error("read should not be permitted")
	}
}

func TestDedupReferenceList(t *testing.T) {
	refs := []acmogo.Reference{
		{Col: "A"},
		{Col: "B"},
		{Col: "C"},
		{Col: "A"},
		{Col: "B"},
		{Col: "D"},
		{Col: "C"},
		{Col: "C"},
		{Col: "D"},
		{Col: "C"},
	}

	refs = acmogo.DedupReferenceList(refs)

	var str string
	for _, ref := range refs {
		str += ref.Col
	}

	if str != "ABDC" {
		t.Errorf("expected %q but got %q", "ABDC", str)
	}
}
