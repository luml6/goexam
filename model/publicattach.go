package model

//公共附件信息结构及单个附件操作

import (
	"errors"
	"reflect"
	"time"

	"notepad-api/utils"

	"gopkg.in/mgo.v2/bson"
)

//附件信息
type PublicAttach struct {
	ObjectId_  bson.ObjectId `bson:"_id"`
	ID         string        `bson:"ID"`         //附件唯一标识
	Key        string        `bson:"key"`        //附件key
	Type       string        `bson:"Type"`       //附件类型： 0-背景图片;1-控件;2-模板
	Status     uint8         `bson:"Status"`     //附件状态: 0-状态正常；1-不可用
	CreateTime time.Time     `bson:"CreateTime"` //创建时间
	UpdateTime time.Time     `bson:"UpdateTime"` //上次更新时间
}

func NewPublicAttach(id string) *PublicAttach {
	att := new(PublicAttach)
	att.ID = id

	return att
}

//添加附件
func (attach *PublicAttach) Add() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	attach.ObjectId_ = bson.NewObjectId()
	attach.CreateTime = utils.NowSecond()
	attach.UpdateTime = utils.NowSecond()

	if err := GetPublicAttachCollection(db).Insert(attach); err != nil {
		return err
	}

	return nil
}

//获取附件信息
func (attach *PublicAttach) Get() (*PublicAttach, error) {
	db := MSessionGet()
	if db == nil {
		return attach, SessionFailError
	}
	defer db.Close()

	if err := GetPublicAttachCollection(db).Find(bson.M{"ID": attach.ID}).One(attach); err != nil {
		if IsMongoNotFound(err) {
			return nil, nil
		}

		return attach, err
	}

	return attach, nil
}

//检查note是否存在
func (attach *PublicAttach) IsExist() bool {

	attTmp := NewPublicAttach(attach.ID)
	attTmp, err := attTmp.Get()

	if err == nil && attTmp != nil {
		return true
	}

	return false
}

func (attach *PublicAttach) updateBson(bsm bson.M) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetPublicAttachCollection(db).Update(bson.M{"ID": attach.ID}, bsm)
}

//更新附件信息
func (attach *PublicAttach) Update(fields []string) error {
	if len(fields) == 0 {
		return errors.New("len of fields is 0")
	}

	var fieldsMap bson.M = make(bson.M, 3)

	for _, v := range fields {
		fieldsMap[v] = reflect.ValueOf(*attach).FieldByName(v).Interface()
	}

	return attach.updateBson(bson.M{"$set": fieldsMap})
}

//彻底删除附件
func (attach *PublicAttach) Remove() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetPublicAttachCollection(db).Remove(bson.M{"ID": attach.ID})
}
