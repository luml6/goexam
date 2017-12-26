package model

//控件信息结构

import (
	"errors"
	"reflect"
	"time"

	"notepad-api/utils"

	"gopkg.in/mgo.v2/bson"
)

//控件信息
type Widget struct {
	ObjectId_  bson.ObjectId `bson:"_id"`
	ID         string        `bson:"ID"`         //附件唯一标识
	Name       string        `bson:"Name"`       //控件名称
	CateName   string        `bson:"CateName"`   //控件分类名称
	Price      int           `bson:"Price"`      //定价
	Status     uint8         `bson:"Status"`     //控件状态: 0-状态正常；1-不可用
	CreateTime time.Time     `bson:"CreateTime"` //创建时间
	UpdateTime time.Time     `bson:"UpdateTime"` //上次更新时间
}

func NewWidget(id string) *Widget {
	att := new(Widget)
	att.ID = id

	return att
}

//添加控件
func (widget *Widget) Add() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	widget.ObjectId_ = bson.NewObjectId()
	widget.CreateTime = utils.NowSecond()

	if err := GetWidgetCollection(db).Insert(widget); err != nil {
		return err
	}

	return nil
}

//获取控件信息
func (widget *Widget) Get() (*Widget, error) {
	db := MSessionGet()
	if db == nil {
		return widget, SessionFailError
	}
	defer db.Close()

	if err := GetWidgetCollection(db).Find(bson.M{"ID": widget.ID, "Status": STATUS_NORMAL}).One(widget); err != nil {
		if IsMongoNotFound(err) {
			return nil, nil
		}

		return widget, err
	}

	return widget, nil
}

//获取控件信息
func (widget *Widget) GetList() ([]Widget, error) {
	var widgetlist []Widget
	db := MSessionGet()
	if db == nil {
		return widgetlist, SessionFailError
	}
	defer db.Close()

	if err := GetWidgetCollection(db).Find(bson.M{"Status": STATUS_NORMAL}).All(&widgetlist); err != nil {
		if IsMongoNotFound(err) {
			return nil, nil
		}

		return widgetlist, err
	}

	return widgetlist, nil
}

//检查控件是否存在
func (widget *Widget) IsExist() bool {

	attTmp := NewWidget(widget.ID)
	attTmp, err := attTmp.Get()

	if err == nil && attTmp != nil {
		return true
	}

	return false
}

func (widget *Widget) updateBson(bsm bson.M) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetWidgetCollection(db).Update(bson.M{"ID": widget.ID}, bsm)
}

//更新控件信息
func (widget *Widget) Update(fields []string) error {
	if len(fields) == 0 {
		return errors.New("len of fields is 0")
	}

	var fieldsMap bson.M = make(bson.M, 3)

	for _, v := range fields {
		fieldsMap[v] = reflect.ValueOf(*widget).FieldByName(v).Interface()
	}

	return widget.updateBson(bson.M{"$set": fieldsMap})
}

//删除信息
func (widget *Widget) Delete() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return widget.updateBson(bson.M{"$set": bson.M{"Status": STATUS_DELETED}})
}

//彻底删除控件
func (widget *Widget) Remove() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetWidgetCollection(db).Remove(bson.M{"ID": widget.ID})
}
