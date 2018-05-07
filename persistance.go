package acmogo

import "github.com/globalsign/mgo"

func InsertList(db *mgo.Database, entityList ...Referencer) (int, error) {
	for i, entity := range entityList {
		ref := entity.Ref()
		if err := db.C(ref.Col).Insert(entity); err != nil {
			return len(entityList) - i, err
		}
	}
	return len(entityList), nil
}

func RefreshEntity(db *mgo.Database, entity Referencer) error {
	ref := entity.Ref()
	return db.C(ref.Col).FindId(ref.ID).One(entity)
}

func UpdateEntity(db *mgo.Database, entity Referencer, updateDoc Map) error {
	ref := entity.Ref()
	return db.C(ref.Col).UpdateId(ref.ID, updateDoc)
}

func ReadPermitted(db *mgo.Database, entity Referencer, refs ...Referencer) bool {
	var (
		ent Entity
		ref = entity.Ref()
	)
	if err := db.C(ref.Col).FindId(ref.ID).Select(SelectEntityDoc).One(&ent); err != nil {
		return false
	}
	return ent.AC.ReadPermitted(refs...)
}

func UpdatePermitted(db *mgo.Database, entity Referencer, refs ...Referencer) bool {
	var (
		ent Entity
		ref = entity.Ref()
	)
	if err := db.C(ref.Col).FindId(ref.ID).Select(SelectEntityDoc).One(&ent); err != nil {
		return false
	}
	return ent.AC.UpdatePermitted(refs...)
}

func DeletePermitted(db *mgo.Database, entity Referencer, refs ...Referencer) bool {
	var (
		ent Entity
		ref = entity.Ref()
	)
	if err := db.C(ref.Col).FindId(ref.ID).Select(SelectEntityDoc).One(&ent); err != nil {
		return false
	}
	return ent.AC.DeletePermitted(refs...)
}

func PersistClearAccessControl(db *mgo.Database, entity Referencer, entities ...Referencer) error {
	ref := entity.Ref()
	var refs []Reference
	for _, ent := range entities {
		refs = append(refs, ent.Ref())
	}
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$pullAll": Map{ACPath + ".u": refs, ACPath + ".d": refs, ACPath + ".r": refs},
	})
}

func PersistPermitRead(db *mgo.Database, entity Referencer, entities ...Referencer) error {
	ref := entity.Ref()
	var refs []Reference
	for _, ent := range entities {
		refs = append(refs, ent.Ref())
	}
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$pullAll":  Map{ACPath + ".u": refs, ACPath + ".d": refs},
		"$addToSet": Map{ACPath + ".r": Map{"$each": refs}},
	})
}

func PersistPermitUpdate(db *mgo.Database, entity Referencer, entities ...Referencer) error {
	ref := entity.Ref()
	var refs []Reference
	for _, ent := range entities {
		refs = append(refs, ent.Ref())
	}
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$pullAll":  Map{ACPath + ".r": refs, ACPath + ".d": refs},
		"$addToSet": Map{ACPath + ".u": Map{"$each": refs}},
	})
}

func PersistPermitDelete(db *mgo.Database, entity Referencer, entities ...Referencer) error {
	ref := entity.Ref()
	var refs []Reference
	for _, ent := range entities {
		refs = append(refs, ent.Ref())
	}
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$pullAll":  Map{ACPath + ".r": refs, ACPath + ".u": refs},
		"$addToSet": Map{ACPath + ".d": Map{"$each": refs}},
	})
}

func PersistPublic(db *mgo.Database, entity Referencer) error {
	ref := entity.Ref()
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$set": Map{ACPath + ".p": true},
	})
}

func PersistPrivate(db *mgo.Database, entity Referencer) error {
	ref := entity.Ref()
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$set": Map{ACPath + ".p": false},
	})
}
