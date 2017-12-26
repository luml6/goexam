package model

import (
	"gopkg.in/mgo.v2/bson"
)

type UserCollectionList struct {
	UID string
}

func NewUserCollectionList(notePath string) *UserCollectionList {
	al := new(UserCollectionList)
	al.UID = notePath

	return al
}

//获取账号下所有registerID
func (al *UserCollectionList) GetList() ([]UserCollection, error) {
	var attList []UserCollection
	db := MSessionGet()
	if db == nil {
		return attList, SessionFailError
	}
	defer db.Close()

	if err := GetUserRegisterCollection(db).Find(bson.M{"UID": al.UID, "Status": STATUS_NORMAL}).All(&attList); err != nil {
		return attList, err
	}

	return attList, nil
}
