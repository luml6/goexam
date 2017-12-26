package model

import (
	"errors"
	"reflect"
	"time"

	"notepad-api/utils"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//用户信息
type User struct {
	ObjectId_    bson.ObjectId `bson:"_id"`
	UID          string        `bson:"UID"` //用户唯一标示
	AccessToken  string        `bson:"AccessToken"`
	RefreshToken string        `bson:"RefreshToken"`
	//	Account      string        `bson:"Account"`
	SyncTime   time.Time `bson:"SyncTime"`   //上次同步时间 2016-09-19 12:00:01
	UpdateTime time.Time `bson:"UpdateTime"` //最近修改时间
	Usn        int       `bson:"Usn"`        //笔记版本信息
}

func NewUser(UID string) *User {
	user := new(User)
	user.UID = UID

	return user
}

//添加用户信息
func (user *User) Add() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	user.ObjectId_ = bson.NewObjectId()
	tn := utils.NowSecond()
	user.UpdateTime = tn
	user.SyncTime = tn

	if err := GetUserCollection(db).Insert(user); err != nil {
		return err
	}

	return nil
}

//检查user是否存在
func (user *User) IsExist() bool {

	userTmp := NewUser(user.UID)
	userTmp, err := userTmp.Get()

	if err == nil && userTmp != nil {
		return true
	}

	return false
}

//获取用户信息
func (user *User) Get() (*User, error) {
	db := MSessionGet()
	if db == nil {
		return user, SessionFailError
	}
	defer db.Close()

	if err := GetUserCollection(db).Find(bson.M{"UID": user.UID}).One(user); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return user, err
	}

	return user, nil
}

func (user *User) updateBson(bsm bson.M) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()
	return GetUserCollection(db).Update(bson.M{"UID": user.UID}, bsm)
}

//更新用户信息
func (user *User) Update(fields []string) error {
	if len(fields) == 0 {
		return errors.New("len of fields is 0")
	}

	var fieldsMap bson.M = make(bson.M, 3)

	for _, v := range fields {
		fieldsMap[v] = reflect.ValueOf(*user).FieldByName(v).Interface()
	}

	return user.updateBson(bson.M{"$set": fieldsMap})
}

//更新用户所有信息
func (user *User) UpdateAll() error {
	var (
		data []byte
		err  error
		bsm  bson.M
	)

	if data, err = bson.Marshal(user); err != nil {
		return err
	}

	if err = bson.Unmarshal(data, &bsm); err != nil {
		return err
	}

	return user.updateBson(bsm)
}
