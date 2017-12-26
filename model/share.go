package model

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"notepad-api/utils"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//分享链接
type Share struct {
	ObjectId_   bson.ObjectId `bson:"_id"`
	ShareID     string        `bson:"ShareID"`     //分享标示
	NoteID      string        `bson:"NoteID"`      //笔记ID
	Question    string        `bson:"Question"`    //趣味密码分享问题
	Answer      string        `bson:"Answer"`      //趣味密码分享答案
	Password    string        `bson:"Password"`    //分享密码
	Type        uint8         `bson:"Type"`        //分享类型
	Status      uint8         `bson:"Status"`      //分享状态
	CreateTime  time.Time     `bson:"CreateTime"`  //上次同步时间 2016-09-19 12:00:01
	FailureTime time.Time     `bson:"FailureTime"` //失效时间单位为分
	OpenNum     int           `bson:"OpenNum"`     //分享链接打开次数
}

func NewShare(shareID string) *Share {
	share := new(Share)
	share.ShareID = shareID

	return share
}

//添加分享信息
func (share *Share) Add() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	share.ObjectId_ = bson.NewObjectId()
	tn := utils.NowSecond()
	share.CreateTime = tn

	if err := GetShareCollection(db).Insert(share); err != nil {
		return err
	}

	return nil
}

//获取分享信息
func (share *Share) Get() (*Share, error) {
	db := MSessionGet()
	if db == nil {
		return share, SessionFailError
	}
	defer db.Close()

	if err := GetShareCollection(db).Find(bson.M{"ShareID": share.ShareID, "Status": STATUS_NORMAL}).One(share); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return share, err
	}

	return share, nil
}

//设置笔记分享状态
func (share *Share) SetStatus(status uint8) (err error) {
	return share.updateBson(bson.M{"$set": bson.M{"Status": status}})
}

//取消笔记分享
func (share *Share) CancleShare() (err error) {
	return share.updateBson(bson.M{"$set": bson.M{"Status": STATUS_DELETED, "Question": ""}})
}

//获取分享信息
func (share *Share) GetByNoteIDAndType() (*Share, error) {
	db := MSessionGet()
	if db == nil {
		return nil, SessionFailError
	}
	defer db.Close()
	var sharetmp *Share
	if err := GetShareCollection(db).Find(bson.M{"NoteID": share.NoteID, "Status": STATUS_NORMAL, "Type": share.Type}).One(&sharetmp); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}
		fmt.Println(sharetmp, err)
		return sharetmp, err
	}

	return sharetmp, nil
}

//获取分享信息
func (share *Share) GetByNoteID() (*Share, error) {
	db := MSessionGet()
	if db == nil {
		return share, SessionFailError
	}
	defer db.Close()

	if err := GetShareCollection(db).Find(bson.M{"NoteID": share.NoteID, "Status": STATUS_NORMAL}).One(share); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return share, err
	}

	return share, nil
}

//获取分享信息
func (share *Share) GetByNoteIDDelete() (*Share, error) {
	db := MSessionGet()
	if db == nil {
		return share, SessionFailError
	}
	defer db.Close()
	shareTmp := NewShare(share.ShareID)
	if err := GetShareCollection(db).Find(bson.M{"NoteID": share.NoteID, "Status": STATUS_DELETED}).One(shareTmp); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return shareTmp, err
	}

	return shareTmp, nil
}

//获取分享信息
func (share *Share) IsExist() bool {
	shareTmp := NewShare(share.ShareID)
	shareTmp, err := shareTmp.Get()

	if err == nil && shareTmp != nil {
		return true
	}

	return false
}

func (share *Share) updateBson(bsm bson.M) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()
	return GetShareCollection(db).Update(bson.M{"ShareID": share.ShareID}, bsm)
}

//更新分享信息
func (share *Share) Update(fields []string) error {
	if len(fields) == 0 {
		return errors.New("len of fields is 0")
	}

	var fieldsMap bson.M = make(bson.M, 3)

	for _, v := range fields {
		fieldsMap[v] = reflect.ValueOf(*share).FieldByName(v).Interface()
	}

	return share.updateBson(bson.M{"$set": fieldsMap})
}
