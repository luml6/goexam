package model

import (
	"gopkg.in/mgo.v2/bson"
)

type OptionList struct {
	VoteID string
}

func NewOptionList(VoteID string) *OptionList {
	al := new(OptionList)
	al.VoteID = VoteID

	return al
}

//获取所有选项
func (al *OptionList) GetList() ([]VoteOption, error) {
	var optionList []VoteOption
	db := MSessionGet()
	if db == nil {
		return optionList, SessionFailError
	}
	defer db.Close()

	if err := GetVoteOptionCollection(db).Find(bson.M{"VoteID": al.VoteID}).All(&optionList); err != nil {

		return optionList, err
	}

	return optionList, nil
}

//删除相同投票下的所有答案
func (al *OptionList) RemoveAll() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	if _, err := GetVoteOptionCollection(db).RemoveAll(bson.M{"VoteID": al.VoteID}); err != nil {

		return err
	}

	return nil
}
