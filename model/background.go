package model

//背景信息结构

import (
	"errors"
	"reflect"
	"time"

	"notepad-api/utils"

	"gopkg.in/mgo.v2/bson"
)

//背景信息
type Background struct {
	ObjectId_  bson.ObjectId `bson:"_id"`
	ID         string        `bson:"ID"`         //附件唯一标识
	Name       string        `bson:"Name"`       //背景名称
	CateName   string        `bson:"CateName"`   //背景分类名称
	Content    string        `bson:"Content"`    //背景内容
	Price      int           `bson:"Price"`      //定价
	Status     uint8         `bson:"Status"`     //附件状态: 0-状态正常；1-不可用
	CreateTime time.Time     `bson:"CreateTime"` //创建时间
	UpdateTime time.Time     `bson:"UpdateTime"` //上次更新时间
}

func NewBackground(id string) *Background {
	att := new(Background)
	att.ID = id

	return att
}

//添加背景
func (paper *Background) Add() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	paper.ObjectId_ = bson.NewObjectId()
	paper.CreateTime = utils.NowSecond()

	if err := GetBackgroundCollection(db).Insert(paper); err != nil {
		return err
	}

	return nil
}

//获取背景信息
func (paper *Background) Get() (*Background, error) {
	db := MSessionGet()
	if db == nil {
		return paper, SessionFailError
	}
	defer db.Close()

	if err := GetBackgroundCollection(db).Find(bson.M{"ID": paper.ID, "Status": STATUS_NORMAL}).One(paper); err != nil {
		if IsMongoNotFound(err) {
			return nil, nil
		}

		return paper, err
	}

	return paper, nil
}

//检查note是否存在
func (paper *Background) IsExist() bool {

	attTmp := NewBackground(paper.ID)
	attTmp, err := attTmp.Get()

	if err == nil && attTmp != nil {
		return true
	}

	return false
}

func (paper *Background) updateBson(bsm bson.M) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetBackgroundCollection(db).Update(bson.M{"ID": paper.ID}, bsm)
}

//获取所有背景信息
func (paper *Background) GetList() ([]Background, error) {
	var temList []Background
	db := MSessionGet()
	if db == nil {
		return temList, SessionFailError
	}
	defer db.Close()

	if err := GetBackgroundCollection(db).Find(bson.M{"Status": STATUS_NORMAL}).All(&temList); err != nil {
		if IsMongoNotFound(err) {
			return nil, nil
		}

		return temList, err
	}

	return temList, nil
}

//更新背景信息
func (paper *Background) Update(fields []string) error {
	if len(fields) == 0 {
		return errors.New("len of fields is 0")
	}

	var fieldsMap bson.M = make(bson.M, 3)

	for _, v := range fields {
		fieldsMap[v] = reflect.ValueOf(*paper).FieldByName(v).Interface()
	}

	return paper.updateBson(bson.M{"$set": fieldsMap})
}

//删除信息
func (paper *Background) Delete() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return paper.updateBson(bson.M{"$set": bson.M{"Status": STATUS_DELETED}})
}

//彻底删除背景
func (paper *Background) Remove() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetBackgroundCollection(db).Remove(bson.M{"ID": paper.ID})
}
