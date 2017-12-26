package model

import (
	"gopkg.in/mgo.v2/bson"
)

type PublicAttachList struct {
	NotePath string
}

func NewPublicAttachList(notePath string) *PublicAttachList {
	al := new(PublicAttachList)
	al.NotePath = notePath

	return al
}

//获取相同笔记下的所有附件
func (al *PublicAttachList) GetList() ([]PublicAttach, error) {
	var attList []PublicAttach
	db := MSessionGet()
	if db == nil {
		return attList, SessionFailError
	}
	defer db.Close()

	if err := GetPublicAttachCollection(db).Find(bson.M{"NotePath": al.NotePath}).All(&attList); err != nil {
		return attList, err
	}

	return attList, nil
}

//删除相同笔记下的所有附件
func (al *PublicAttachList) RemoveAll() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	if _, err := GetPublicAttachCollection(db).RemoveAll(bson.M{"NotePath": al.NotePath}); err != nil {

		return err
	}

	return nil
}
