package model

//附件信息结构及单个附件操作

import (
	"errors"
	"reflect"
	"time"

	"notepad-api/utils"

	"gopkg.in/mgo.v2/bson"
)

//附件信息
type Attachment struct {
	ObjectId_      bson.ObjectId `bson:"_id"`
	ID             string        `bson:"ID"`             //附件唯一路径
	NotePath       string        `bson:"NotePath"`       //附件所属笔记路径
	OK             string        `bson:"OK"`             //object key 附件key,若OT=0，则OK为第三方云的key；若OT=1；则OK为附件URL
	OT             uint8         `bson:"OT"`             //object type 附件存储类型： 0-第三方云；1-服务器
	Name           string        `bson:"Name"`           //附件文件名称
	Type           string        `bson:"Type"`           //附件文件类型
	Status         uint8         `bson:"Status"`         //附件状态
	UID            string        `bson:"UID"`            //附件文件类型
	Source         uint8         `bson:"Source"`         //附件来源
	Size           int           `bson:"Size"`           //附件文件大小
	FileEncryption string        `bson:"FileEncryption"` //附件加密key
	CreateTime     time.Time     `bson:"CreateTime"`     //附件创建时间
}

func NewAttachment(id string) *Attachment {
	att := new(Attachment)
	att.ID = id

	return att
}

//添加附件
func (attach *Attachment) Add() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	attach.ObjectId_ = bson.NewObjectId()
	attach.CreateTime = utils.NowSecond()

	if err := GetAttachCollection(db).Insert(attach); err != nil {
		return err
	}

	return nil
}

//获取附件信息
func (attach *Attachment) Get() (*Attachment, error) {
	db := MSessionGet()
	if db == nil {
		return attach, SessionFailError
	}
	defer db.Close()

	if err := GetAttachCollection(db).Find(bson.M{"ID": attach.ID, "Status": attach.Status}).One(attach); err != nil {
		if IsMongoNotFound(err) {
			return nil, nil
		}

		return attach, err
	}

	return attach, nil
}

//获取附件信息
func (attach *Attachment) GetAllStatus() (*Attachment, error) {
	db := MSessionGet()
	if db == nil {
		return attach, SessionFailError
	}
	defer db.Close()

	if err := GetAttachCollection(db).Find(bson.M{"ID": attach.ID}).One(attach); err != nil {
		if IsMongoNotFound(err) {
			return nil, nil
		}

		return attach, err
	}

	return attach, nil
}

//检查附件是否存在
func (attach *Attachment) IsExist() bool {

	attTmp := NewAttachment(attach.ID)
	attTmp, err := attTmp.GetAllStatus()

	if err == nil && attTmp != nil {
		return true
	}

	return false
}

//设置附件状态
func (attach *Attachment) SetStatus(status uint8) (err error) {
	return attach.updateBson(bson.M{"$set": bson.M{"Status": status}})
}

func (attach *Attachment) updateBson(bsm bson.M) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetAttachCollection(db).Update(bson.M{"ID": attach.ID}, bsm)
}

//更新附件信息
func (attach *Attachment) Update(fields []string) error {
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
func (attach *Attachment) Remove() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetAttachCollection(db).Remove(bson.M{"ID": attach.ID, "NotePath": attach.NotePath})
}
