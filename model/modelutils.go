package model

import (
	"errors"

	"gopkg.in/mgo.v2"
)

var (
	SessionFailError error = errors.New("get session fail")
)

const (
	STATUS_NORMAL = iota
	STATUS_DELETED
	STATUS_PRECREATE
	STATUS_TRASH
)
const (
	NOTE_NOT_SHARE = 0
	NOTE_SHARE     = 1
)

const (
	NORMAL_SHARE   = 1
	PASSWORD_SHARE = 2
	QUESTION_SHARE = 4
	VOTE_SHARE     = 8
)

const (
	NOTE_NOT_PINNED = 0
	NOTE_PINNED     = 1
)

const FOTMAT_TIME_STRING = "2006-01-02 15:04:05"

const (
	ATT_STORE_TYPE_COULD = iota
	ATT_STORE_TYPE_ZEUSIS
)

func GetUserCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("user")
}

func GetCategoryCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("category")
}

func GetNoteCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("note")
}

func GetUserRegisterCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("register")
}

func GetAttachCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("attachment")
}
func GetConflictCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("conflict")
}
func GetBackgroundCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("background")
}

func GetTemplateCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("template")
}

func GetWidgetCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("widget")
}
func GetShareCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("share")
}

func GetPublicAttachCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("publicAttach")
}

func GetShareAnswerCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("shareAnswer")
}

func GetVoteCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("vote")
}

func GetVoteOptionCollection(db *mgo.Session) *mgo.Collection {
	if db == nil {
		return nil
	}

	return db.DB("cloudnote").C("voteOption")
}

func IsMongoNotFound(err error) bool {
	if err.Error() == mgo.ErrNotFound.Error() {
		return true
	}

	return false
}
