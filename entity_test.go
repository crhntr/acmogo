package entity_test

import (
	"flag"
	"os"
	"testing"

	"github.com/crhntr/litsphere/internal/entity"
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
	identityZeroVal := entity.Reference{}
	if identityZeroVal.Validate() == nil {
		t.Fatal()
	}

	user0 := User{Entity: entity.New()}
	user1 := User{Entity: entity.New()}
	user2 := User{Entity: entity.New()}

	team0 := Team{Entity: entity.New()}
	team1 := Team{Entity: entity.New()}

	for _, idfr := range []entity.Referencer{user0, user1, user2, team0, team1} {
		if err := idfr.Ref().Validate(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestInsertEntity(t *testing.T) {
	user0 := User{Entity: entity.New()}
	user1 := User{Entity: entity.New()}
	team0 := Team{Entity: entity.New()}
	post0 := Post{Entity: entity.New()}

	if _, err := entity.InsertList(db, user0, user1, team0, post0); err != nil {
		t.Fatal(err)
	}
	if _, err := entity.InsertList(db, user0); err == nil {
		t.Fatal(err)
	}
}

func TestRefreshEntity(t *testing.T) {
	post0 := Post{Entity: entity.New()}
	entity.InsertList(db, post0)
	post0N := 101
	if err := db.C(PostCol).UpdateId(post0.ID, mp{"$set": mp{"n": post0N}}); err != nil {
		t.Fatal(err)
	}
	// t.Log(post0)
	entity.RefreshEntity(db, &post0)
	// t.Log(post0)
	if post0.N != post0N {
		t.Fail()
	}
}

func TestUpdateEntity(t *testing.T) {
	post0 := Post{Entity: entity.New()}
	entity.InsertList(db, post0)
	post0N := 101
	// t.Log(post0)
	entity.UpdateEntity(db, post0, bson.M{"$set": bson.M{"n": post0N}})
	entity.RefreshEntity(db, &post0)
	// t.Log(post0)
	if post0.N != post0N {
		t.Fail()
	}
}

func TestPersistClearAccessControl(t *testing.T) {
	post0 := Post{Entity: entity.New()}
	user0 := User{Entity: entity.New()}

	post0.PermitRead(user0)
	entity.InsertList(db, post0, user0)

	if err := entity.PersistClearAccessControl(db, post0, user0); err != nil {
		t.Fail()
	}
	if entity.DeletePermitted(db, post0, user0) {
		t.Fail()
	}
	if entity.UpdatePermitted(db, post0, user0) {
		t.Fail()
	}
	if entity.ReadPermitted(db, post0, user0) {
		t.Fail()
	}
}

func TestPersistPermitRead(t *testing.T) {
	post0 := Post{Entity: entity.New()}
	post1 := Post{Entity: entity.New()}
	user0 := User{Entity: entity.New()}
	user1 := User{Entity: entity.New()}

	entity.InsertList(db, post0, user0, user1)

	if entity.ReadPermitted(db, post0, user0, user1) {
		t.Error("read should not be permitted")
	}

	entity.PersistPermitRead(db, post0, user0)

	if !entity.ReadPermitted(db, post0, user0) {
		t.Error("read should be permitted")
	}
	if entity.UpdatePermitted(db, post0, user0) {
		t.Error("update should not be permitted")
	}
	if entity.DeletePermitted(db, post0, user0) {
		t.Error("delete should not be permitted")
	}
	if entity.ReadPermitted(db, post0, user1) {
		t.Error("read should not be permitted for unknown reader")
	}
	if entity.ReadPermitted(db, post1, user1) {
		t.Error("post1 should not be found")
	}
}

func TestPersistPermitUpdate(t *testing.T) {
	post0 := Post{Entity: entity.New()}
	post1 := Post{Entity: entity.New()}
	user0 := User{Entity: entity.New()}
	user1 := User{Entity: entity.New()}

	entity.InsertList(db, post0, user0, user1)

	if entity.UpdatePermitted(db, post0, user0, user1) {
		t.Error("read should not be permitted")
	}

	entity.PersistPermitUpdate(db, post0, user0)

	if !entity.ReadPermitted(db, post0, user0) {
		t.Error("read should be permitted")
	}
	if !entity.UpdatePermitted(db, post0, user0) {
		t.Error("update should be permitted")
	}
	if entity.DeletePermitted(db, post0, user0) {
		t.Error("delete should not be permitted")
	}

	if entity.UpdatePermitted(db, post0, user1) {
		t.Error("update should not be permitted for unknown updater")
	}
	if entity.UpdatePermitted(db, post1, user1) {
		t.Error("post1 should not be found")
	}
}

func TestPersistPermitDelete(t *testing.T) {
	post0 := Post{Entity: entity.New()}
	post1 := Post{Entity: entity.New()}
	user0 := User{Entity: entity.New()}
	user1 := User{Entity: entity.New()}

	entity.InsertList(db, post0, user0, user1)

	if entity.DeletePermitted(db, post0, user0, user1) {
		t.Error("delete should not be permitted")
	}

	entity.PersistPermitDelete(db, post0, user0)

	if !entity.ReadPermitted(db, post0, user0) {
		t.Error("read should be permitted")
	}
	if !entity.UpdatePermitted(db, post0, user0) {
		t.Error("update should be permitted")
	}
	if !entity.DeletePermitted(db, post0, user0) {
		t.Error("delete should be permitted")
	}

	if entity.DeletePermitted(db, post0, user1) {
		t.Error("delete should not be permitted for unknown deleter")
	}
	if entity.DeletePermitted(db, post1, user1) {
		t.Fail()
	}
}

func TestPersistPermitDeleteDowngradeToPermitUpdate(t *testing.T) {
	post0 := Post{Entity: entity.New()}
	user0 := User{Entity: entity.New()}

	entity.InsertList(db, post0, user0)

	entity.PersistPermitDelete(db, post0, user0)
	entity.PersistPermitUpdate(db, post0, user0)

	if !entity.ReadPermitted(db, post0, user0) {
		t.Fail()
	}
	if !entity.UpdatePermitted(db, post0, user0) {
		t.Fail()
	}
	if entity.DeletePermitted(db, post0, user0) {
		t.Fail()
	}
}

func TestPersistPermitDeleteDowngradeToPermitRead(t *testing.T) {
	post0 := Post{Entity: entity.New()}
	user0 := User{Entity: entity.New()}

	entity.InsertList(db, post0, user0)

	entity.PersistPermitDelete(db, post0, user0)
	entity.PersistPermitRead(db, post0, user0)

	if !entity.ReadPermitted(db, post0, user0) {
		t.Fail()
	}
	if entity.UpdatePermitted(db, post0, user0) {
		t.Fail()
	}
	if entity.DeletePermitted(db, post0, user0) {
		t.Fail()
	}
}

func TestPersistPermitUpdateDowngradeToPermitRead(t *testing.T) {
	post0 := Post{Entity: entity.New()}
	user0 := User{Entity: entity.New()}

	entity.InsertList(db, post0, user0)

	entity.PersistPermitUpdate(db, post0, user0)
	entity.PersistPermitRead(db, post0, user0)

	if !entity.ReadPermitted(db, post0, user0) {
		t.Fail()
	}
	if entity.UpdatePermitted(db, post0, user0) {
		t.Fail()
	}
	if entity.DeletePermitted(db, post0, user0) {
		t.Fail()
	}
}

func TestPersistPublic(t *testing.T) {
	post0 := Post{Entity: entity.New()}
	user0 := User{Entity: entity.New()}
	user1 := User{Entity: entity.New()}

	entity.InsertList(db, post0, user0, user1)

	if entity.ReadPermitted(db, post0, user0, user1) {
		t.Error("read should not be permitted")
	}

	entity.PersistPublic(db, post0)

	if entity.ReadPermitted(db, post0, user0, user1) {
		t.Error("read should be permitted")
	}
}

func TestPersistPrivate(t *testing.T) {
	post0 := Post{Entity: entity.New()}
	user0 := User{Entity: entity.New()}
	user1 := User{Entity: entity.New()}

	entity.InsertList(db, post0, user0, user1)

	if entity.ReadPermitted(db, post0, user0, user1) {
		t.Error("read should be permitted")
	}

	entity.PersistPrivate(db, post0)

	if entity.ReadPermitted(db, post0, user0, user1) {
		t.Error("read should not be permitted")
	}
}

func TestDedupReferenceList(t *testing.T) {
	refs := []entity.Reference{
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

	refs = entity.DedupReferenceList(refs)

	var str string
	for _, ref := range refs {
		str += ref.Col
	}

	if str != "ABDC" {
		t.Errorf("expected %q but got %q", "ABDC", str)
	}
}
