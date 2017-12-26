package model

import (
	"errors"
	"reflect"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//分享链接
type VoteOption struct {
	ObjectId_   bson.ObjectId `bson:"_id"`
	ID          string        `bson:"ID"`          //答案ID
	VoteID      string        `bson:"VoteID"`      //分享标示
	Option      string        `bson:"Option"`      //投票选项
	OptionCount int           `bson:"OptionCount"` //投票次数
}

func NewVoteOption(ID string) *VoteOption {
	option := new(VoteOption)
	option.ID = ID

	return option
}

//添加分享信息
func (option *VoteOption) Add() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	option.ObjectId_ = bson.NewObjectId()

	if err := GetVoteOptionCollection(db).Insert(option); err != nil {
		return err
	}

	return nil
}

//获取投票信息
func (option *VoteOption) IsExist() bool {
	optionTmp := NewVoteOption(option.ID)
	optionTmp, err := optionTmp.Get()

	if err == nil && optionTmp != nil {
		return true
	}

	return false
}

//获取分享信息
func (option *VoteOption) Get() (*VoteOption, error) {
	db := MSessionGet()
	if db == nil {
		return option, SessionFailError
	}
	defer db.Close()

	if err := GetVoteOptionCollection(db).Find(bson.M{"ID": option.ID}).One(option); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return option, err
	}

	return option, nil
}

func (option *VoteOption) updateBson(bsm bson.M) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()
	return GetVoteOptionCollection(db).Update(bson.M{"ID": option.ID}, bsm)
}

//更新分享信息
func (option *VoteOption) Update(fields []string) error {
	if len(fields) == 0 {
		return errors.New("len of fields is 0")
	}

	var fieldsMap bson.M = make(bson.M, 3)

	for _, v := range fields {
		fieldsMap[v] = reflect.ValueOf(*option).FieldByName(v).Interface()
	}

	return option.updateBson(bson.M{"$set": fieldsMap})
}
