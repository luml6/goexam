package model

import (
	"gopkg.in/mgo.v2/bson"
)

type AttachmentList struct {
	NotePath string
}

func NewAttachmentList(notePath string) *AttachmentList {
	al := new(AttachmentList)
	al.NotePath = notePath

	return al
}

//获取相同笔记下的所有附件
func (al *AttachmentList) GetList() ([]Attachment, error) {
	var attList []Attachment
	db := MSessionGet()
	if db == nil {
		return attList, SessionFailError
	}
	defer db.Close()

	if err := GetAttachCollection(db).Find(bson.M{"NotePath": al.NotePath, "Status": STATUS_NORMAL}).All(&attList); err != nil {
		return attList, err
	}

	return attList, nil
}

//删除相同笔记下的所有附件
func (al *AttachmentList) RemoveAll() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	if _, err := GetAttachCollection(db).RemoveAll(bson.M{"NotePath": al.NotePath}); err != nil {

		return err
	}

	return nil
}
