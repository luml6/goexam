package model

import (
	"notepad-api/utils"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Conflict struct {
	ObjectId_  bson.ObjectId `bson:"_id"`
	ID         string        `bson:"ID"`
	FatherPath string        `bson:"fatherPath"`
	SonPath    string        `bson:"sonPath"`
	CheckSum   string        `bson:"checkSum"`
	CreateTime time.Time     `bson:"createTime"`
}

func NewConflict(id string) *Conflict {
	att := new(Conflict)
	att.ID = id

	return att
}

//添加冲突处理
func (attach *Conflict) Add() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	attach.ObjectId_ = bson.NewObjectId()
	attach.CreateTime = utils.NowSecond()

	if err := GetConflictCollection(db).Insert(attach); err != nil {
		return err
	}

	return nil
}

func (con *Conflict) GetBySum() (*Conflict, error) {
	conTmp := new(Conflict)
	db := MSessionGet()
	if db == nil {
		return conTmp, SessionFailError
	}
	defer db.Close()

	if err := GetConflictCollection(db).Find(bson.M{"checkSum": con.CheckSum, "fatherPath": con.FatherPath}).One(conTmp); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return conTmp, err
	}
	return conTmp, nil
}
