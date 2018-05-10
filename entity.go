package acmogo

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
)

type Entity struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	AC        `json:"_ac" bson:"_ac"`
	CreatedAt time.Time `json:"_createdAt" bson:"_createdAt"`
	// UpdatedAt time.Time `json:"_updatedAt" bson:"_updatedAt"`
}

func New() Entity {
	return Entity{
		ID:        bson.NewObjectId(),
		CreatedAt: time.Now(),
	}
}

type Reference struct {
	Col string        `json:"c" bson:"c"`
	ID  bson.ObjectId `json:"id" bson:"id"`
}

type Referencer interface {
	Ref() Reference
}

type Map = bson.M

var SelectEntityDoc = map[string]int{"_id": 1, ACPath: 1}

func (ref Reference) Validate() error {
	if ref.Col == "" || !ref.ID.Valid() {
		return fmt.Errorf("invalid identity {%q: %q}", ref.Col, ref.ID)
	}
	return nil
}

func (ref Reference) Ref() Reference {
	return ref
}

func FilterReferenceList(ids []Reference, cutset ...Reference) []Reference {
	filtered := ids[:0]
	for _, id := range cutset {
		for _, idx := range ids {
			if id.ID != idx.ID || id.Col != idx.Col {
				filtered = append(filtered, idx)
			}
		}
	}
	return filtered
}

// DedupReferenceList takes a slice of Reference and returns
// a deduplicated slice using a new underlying array.
func DedupReferenceList(refs []Reference) []Reference {
	// this would use the same underlying array.
	// Which may be an unexpected side affect
	// results := refs[:0]
	// instead we should allocate a new underlying array
	// and pesimally allocated an underlyng array of the same size as refs
	results := make([]Reference, 0, len(refs))

	for i, ref := range refs {
		exists := false
		for _, pref := range refs[i+1:] {
			if ref.Col == pref.Col && ref.ID == pref.ID {
				exists = true
				break
			}
		}
		if !exists {
			results = append(results, ref)
		}
	}
	return results
}

func MakeReferenceList(col string, ids ...bson.ObjectId) []Reference {
	refs := make([]Reference, len(ids))
	for i, id := range ids {
		refs[i] = Reference{col, id}
	}
	return refs
}
