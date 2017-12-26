package model

//模板信息结构

import (
	"errors"
	"reflect"
	"time"

	"notepad-api/utils"

	"gopkg.in/mgo.v2/bson"
)

//模板信息
type Template struct {
	ObjectId_  bson.ObjectId `bson:"_id"`
	ID         string        `bson:"ID"`         //附件唯一标识
	Name       string        `bson:"Name"`       //模板名称
	CateName   string        `bson:"CateName"`   //模板分类名称
	Content    string        `bson:"Content"`    //模板内容
	Price      int           `bson:"Price"`      //定价
	Status     uint8         `bson:"Status"`     //附件状态: 0-状态正常；1-不可用
	CreateTime time.Time     `bson:"CreateTime"` //创建时间
	UpdateTime time.Time     `bson:"UpdateTime"` //上次更新时间
}

func NewTemplate(id string) *Template {
	att := new(Template)
	att.ID = id

	return att
}

//添加模板
func (template *Template) Add() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	template.ObjectId_ = bson.NewObjectId()
	template.CreateTime = utils.NowSecond()
	template.UpdateTime = utils.NowSecond()
	template.Status = STATUS_NORMAL

	if err := GetTemplateCollection(db).Insert(template); err != nil {
		return err
	}

	return nil
}

//获取模板信息
func (template *Template) Get() (*Template, error) {
	db := MSessionGet()
	if db == nil {
		return template, SessionFailError
	}
	defer db.Close()

	if err := GetTemplateCollection(db).Find(bson.M{"ID": template.ID, "Status": STATUS_NORMAL}).One(template); err != nil {
		if IsMongoNotFound(err) {
			return nil, nil
		}

		return template, err
	}

	return template, nil
}

//获取所有模板信息
func (template *Template) GetList() ([]Template, error) {
	var temList []Template
	db := MSessionGet()
	if db == nil {
		return temList, SessionFailError
	}
	defer db.Close()

	if err := GetTemplateCollection(db).Find(bson.M{"Status": STATUS_NORMAL}).All(&temList); err != nil {

		if IsMongoNotFound(err) {
			return nil, nil
		}

		return temList, err
	}

	return temList, nil
}

//检查note是否存在
func (template *Template) IsExist() bool {

	attTmp := NewTemplate(template.ID)
	attTmp, err := attTmp.Get()

	if err == nil && attTmp != nil {
		return true
	}

	return false
}

func (template *Template) updateBson(bsm bson.M) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetTemplateCollection(db).Update(bson.M{"ID": template.ID}, bsm)
}

//更新附件信息
func (template *Template) Update(fields []string) error {
	if len(fields) == 0 {
		return errors.New("len of fields is 0")
	}

	var fieldsMap bson.M = make(bson.M, 3)

	for _, v := range fields {
		fieldsMap[v] = reflect.ValueOf(*template).FieldByName(v).Interface()
	}

	return template.updateBson(bson.M{"$set": fieldsMap})
}

//删除附件信息
func (template *Template) Delete() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return template.updateBson(bson.M{"$set": bson.M{"Status": STATUS_DELETED}})
}

//彻底删除附件
func (template *Template) Remove() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetTemplateCollection(db).Remove(bson.M{"ID": template.ID})
}
