package acmogo

import (
	"errors"
)

var (
	ACPath       = "_ac"
	PublicPath   = ACPath + ".pu"
	ReadersPath  = ACPath + ".r"
	UpdatersPath = ACPath + ".u"
	DeletersPath = ACPath + ".d"
	CreatorPath  = ACPath + ".cr"
)

// AC should be Embeded in structs to be stored in MongoDB
// It should be anotated with the `bson:"ac"` or whatever ACPath is set to.
// When a new object is created, the creator's identity should be passed to SetCreator
// bson tag "inline" should not be set
type AC struct {
	Readers  []Reference `json:"r,omitempty" bson:"r,omitempty"`
	Updaters []Reference `json:"u,omitempty" bson:"u,omitempty"`
	Deleters []Reference `json:"d,omitempty" bson:"d,omitempty"`
	Creator  *Reference  `json:"cr,omitempty" bson:"cr,omitempty"`
	Public   bool        `json:"pu" bson:"pu"`
}

func (ac *AC) SetCreator(id Reference) error {
	if ac.Creator != nil {
		return errors.New("creator already set")
	}
	if err := id.Validate(); err != nil {
		return err
	}
	ac.Creator = &id
	return nil
}

func (ac AC) ReadPermitted(refs ...Referencer) bool {
	if ac.Public {
		return true
	}
	for _, ref := range refs {
		id := ref.Ref()
		if ac.Creator != nil && ac.Creator.Col == id.Col && ac.Creator.ID == id.ID {
			return true
		}

		for _, idInSet := range ac.Deleters {
			if idInSet.Col == id.Col && idInSet.ID == id.ID {
				return true
			}
		}
		for _, idInSet := range ac.Updaters {
			if idInSet.Col == id.Col && idInSet.ID == id.ID {
				return true
			}
		}
		for _, idInSet := range ac.Readers {
			if idInSet.Col == id.Col && idInSet.ID == id.ID {
				return true
			}
		}
	}
	return false
}

func (ac AC) UpdatePermitted(refs ...Referencer) bool {
	for _, ref := range refs {
		id := ref.Ref()
		if ac.Creator != nil && ac.Creator.Col == id.Col && ac.Creator.ID == id.ID {
			return true
		}

		for _, idInSet := range ac.Deleters {
			if idInSet.Col == id.Col && idInSet.ID == id.ID {
				return true
			}
		}
		for _, idInSet := range ac.Updaters {
			if idInSet.Col == id.Col && idInSet.ID == id.ID {
				return true
			}
		}
	}
	return false
}

func (ac AC) DeletePermitted(refs ...Referencer) bool {
	for _, ref := range refs {
		id := ref.Ref()
		if ac.Creator != nil && ac.Creator.Col == id.Col && ac.Creator.ID == id.ID {
			return true
		}

		for _, idInSet := range ac.Deleters {
			if idInSet.Col == id.Col && idInSet.ID == id.ID {
				return true
			}
		}
	}
	return false
}

func (ac *AC) ClearAccessControl(refs ...Referencer) {
	for _, ref := range refs {
		r := ref.Ref()
		ac.Updaters = FilterReferenceList(ac.Updaters, r)
		ac.Deleters = FilterReferenceList(ac.Deleters, r)
		ac.Readers = FilterReferenceList(ac.Readers, r)
	}
}

func (ac *AC) PermitRead(refs ...Referencer) {
	for _, ref := range refs {
		r := ref.Ref()
		ac.Updaters = FilterReferenceList(ac.Updaters, r)
		ac.Deleters = FilterReferenceList(ac.Deleters, r)
		ac.Readers = FilterReferenceList(ac.Readers, r)

		ac.Readers = append(ac.Readers, r)
	}
}

func (ac *AC) PermitUpdate(refs ...Referencer) {
	for _, ref := range refs {
		r := ref.Ref()
		ac.Updaters = FilterReferenceList(ac.Updaters, r)
		ac.Deleters = FilterReferenceList(ac.Deleters, r)
		ac.Readers = FilterReferenceList(ac.Readers, r)

		ac.Updaters = append(ac.Updaters, r)
	}
}

func (ac *AC) PermitDelete(refs ...Referencer) {
	for _, ref := range refs {
		r := ref.Ref()
		ac.Updaters = FilterReferenceList(ac.Updaters, r)
		ac.Deleters = FilterReferenceList(ac.Deleters, r)
		ac.Readers = FilterReferenceList(ac.Readers, r)

		ac.Deleters = append(ac.Deleters, r)
	}
}
