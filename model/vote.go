package model

import (
	"errors"
	"reflect"
	"time"

	"notepad-api/utils"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//投票
type Vote struct {
	ObjectId_   bson.ObjectId `bson:"_id"`
	VoteID      string        `bson:"VoteID"`      //分享标示
	NoteID      string        `bson:"NoteID"`      //笔记ID
	Question    string        `bson:"Question"`    //分享问题
	Type        uint8         `bson:"Type"`        //0表示单选，1表示多选
	CreateTime  time.Time     `bson:"CreateTime"`  //上次同步时间 2016-09-19 12:00:01
	InvalidTime time.Time     `bson:"InvalidTime"` //失效时间
	ClickNum    int           `bson:"ClickNum"`    //投票人数
}

func NewVote(voteID string) *Vote {
	vote := new(Vote)
	vote.VoteID = voteID

	return vote
}

//添加投票信息
func (vote *Vote) Add() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	vote.ObjectId_ = bson.NewObjectId()
	tn := utils.NowSecond()
	vote.CreateTime = tn

	if err := GetVoteCollection(db).Insert(vote); err != nil {
		return err
	}

	return nil
}

//获取投票信息
func (vote *Vote) Get() (*Vote, error) {
	db := MSessionGet()
	if db == nil {
		return vote, SessionFailError
	}
	defer db.Close()
	beego.Debug(vote.VoteID)
	if err := GetVoteCollection(db).Find(bson.M{"VoteID": vote.VoteID}).One(vote); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return vote, err
	}
	return vote, nil
}

//获取投票信息
func (vote *Vote) GetByNoteID() (*Vote, error) {
	db := MSessionGet()
	if db == nil {
		return vote, SessionFailError
	}
	defer db.Close()

	if err := GetVoteCollection(db).Find(bson.M{"NoteID": vote.NoteID}).One(vote); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return vote, err
	}

	return vote, nil
}

//获取投票信息
func (vote *Vote) IsExist() bool {
	voteTmp := NewVote(vote.VoteID)
	voteTmp, err := voteTmp.Get()
	if err == nil && voteTmp != nil {
		return true
	}

	return false
}

func (vote *Vote) updateBson(bsm bson.M) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()
	return GetVoteCollection(db).Update(bson.M{"VoteID": vote.VoteID}, bsm)
}

//删除投票
func (al *Vote) Remove() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	if err := GetVoteCollection(db).Remove(bson.M{"VoteID": al.VoteID}); err != nil {

		return err
	}

	return nil
}

//更新投票信息
func (vote *Vote) Update(fields []string) error {
	if len(fields) == 0 {
		return errors.New("len of fields is 0")
	}

	var fieldsMap bson.M = make(bson.M, 3)

	for _, v := range fields {
		fieldsMap[v] = reflect.ValueOf(*vote).FieldByName(v).Interface()
	}

	return vote.updateBson(bson.M{"$set": fieldsMap})
}
