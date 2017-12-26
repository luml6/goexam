package model

import (
	"errors"
	"notepad-api/utils"
	"reflect"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserCollection struct {
	ObjectId_  bson.ObjectId `bson:"_id"`
	UID        string        `bson:"UID"`
	RegisterId string        `bson:"RegisterId"`
	CreateTime time.Time     `bson:"CreateTime"`
	UpdateTime time.Time     `bson:"UpdateTime"`
	Status     uint8         `bson:"Status"`
}

func NewUserCollection(uid string) *UserCollection {
	user := new(UserCollection)
	user.UID = uid

	return user
}

//添加用户信息
func (user *UserCollection) Add() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	user.ObjectId_ = bson.NewObjectId()
	tn := utils.NowSecond()
	user.UpdateTime = tn
	user.CreateTime = tn
	if err := GetUserRegisterCollection(db).Insert(user); err != nil {
		return err
	}

	return nil
}

//检查user是否存在
func (user *UserCollection) IsExist() bool {

	userTmp := NewUserCollection(user.UID)
	userTmp.RegisterId = user.RegisterId
	userTmp, err := userTmp.Get()

	if err == nil && userTmp != nil {
		return true
	}

	return false
}

//检查user是否存在
func (user *UserCollection) GetByRegisterID() (*UserCollection, error) {
	userTmp := new(UserCollection)
	db := MSessionGet()
	if db == nil {
		return userTmp, SessionFailError
	}
	defer db.Close()

	if err := GetUserRegisterCollection(db).Find(bson.M{"RegisterId": user.RegisterId, "Status": STATUS_NORMAL}).One(userTmp); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return userTmp, err
	}
	return userTmp, nil
}

//获取用户信息
func (user *UserCollection) Get() (*UserCollection, error) {
	db := MSessionGet()
	if db == nil {
		return user, SessionFailError
	}
	defer db.Close()

	if err := GetUserRegisterCollection(db).Find(bson.M{"UID": user.UID, "RegisterId": user.RegisterId, "Status": STATUS_NORMAL}).One(user); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return user, err
	}

	return user, nil
}

func (user *UserCollection) updateBson(bsm bson.M) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetUserRegisterCollection(db).Update(bson.M{"UID": user.UID}, bsm)
}

//更新用户信息
func (user *UserCollection) Update(fields []string) error {
	if len(fields) == 0 {
		return errors.New("len of fields is 0")
	}

	var fieldsMap bson.M = make(bson.M, 3)

	for _, v := range fields {
		fieldsMap[v] = reflect.ValueOf(*user).FieldByName(v).Interface()
	}

	return user.updateBson(bson.M{"$set": fieldsMap})
}
