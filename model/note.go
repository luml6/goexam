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

//笔记信息
type Note struct {
	ObjectId_     bson.ObjectId `bson:"_id"`
	Path          string        `bson:"Path"`          //笔记
	CategoryPath  string        `bson:"CategoryPath"`  //笔记分组唯一途径
	CategoryName  string        `bson:"CategoryName"`  //笔记分组唯一名称
	UID           string        `bson:"UID"`           //用户唯一ID
	Title         string        `bson:"Title"`         //笔记标题
	Author        string        `bson:"Author"`        //笔记作者
	Summary       string        `bson:"Summary"`       //笔记摘要
	Content       string        `bson:"Content"`       //笔记内容
	CheckSum      string        `bson:"CheckSum"`      //笔记校验和
	Source        uint8         `bson:"Source"`        //笔记来源
	Status        uint8         `bson:"Status"`        //笔记状态
	IsPinned      uint8         `bson:"IsPinned"`      //笔记是否置顶
	IsShare       uint8         `bson:"IsShare"`       //笔记是否被分享
	IsPassword    uint8         `bson:"IsPassword"`    //笔记是否被加密分享
	ConflictCount int           `bson:"ConflictCount"` //笔记冲突次数
	IsVote        uint8         `bson:"IsVote"`        //笔记是否投票
	CreateTime    time.Time     `bson:"CreateTime"`    //笔记创建时间
	UpdateTime    time.Time     `bson:"UpdateTime"`    //笔记更新时间

	Usn int `bson:"Usn"` //笔记更新时间
}

func NewNote(path string) *Note {
	note := new(Note)
	note.Path = path

	return note
}

//根据状态获取笔记信息
func (note *Note) GetByPathAndStatus(Status uint8) (*Note, error) {
	db := MSessionGet()
	if db == nil {
		return note, SessionFailError
	}
	defer db.Close()

	if err := GetNoteCollection(db).Find(bson.M{"UID": note.UID, "Path": note.Path, "Status": Status}).One(note); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return note, err
	}

	return note, nil
}

//根据状态获取笔记信息
func (note *Note) GetByPathAndUsn(usn int) (*Note, error) {
	db := MSessionGet()
	if db == nil {
		return note, SessionFailError
	}
	defer db.Close()

	if err := GetNoteCollection(db).Find(bson.M{"UID": note.UID, "Path": note.Path, "Usn": usn}).One(note); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return note, err
	}

	return note, nil
}

//查询笔记
func (note *Note) Get() (*Note, error) {
	return note.GetByPathAndStatus(STATUS_NORMAL)
}

//查询笔记
func (note *Note) GetWithPath() (*Note, error) {
	db := MSessionGet()
	if db == nil {
		return note, SessionFailError
	}
	defer db.Close()

	if err := GetNoteCollection(db).Find(bson.M{"Path": note.Path, "Status": STATUS_NORMAL}).One(note); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return note, err
	}

	return note, nil
}

//查询所有状态笔记
func (note *Note) GetAllStatus() (*Note, error) {
	db := MSessionGet()
	if db == nil {
		return note, SessionFailError
	}
	defer db.Close()

	if err := GetNoteCollection(db).Find(bson.M{"UID": note.UID, "Path": note.Path}).One(note); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			beego.Debug(err)
			return nil, nil
		}

		return note, err
	}

	return note, nil
}

//添加笔记
func (note *Note) Add() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	note.ObjectId_ = bson.NewObjectId()
	//	tn := utils.NowSecond()

	//	note.CreateTime = tn
	//	note.UpdateTime = tn
	note.Status = STATUS_NORMAL

	if err := GetNoteCollection(db).Insert(note); err != nil {
		return err
	}

	return nil
}

//添加笔记
func (note *Note) Create() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	note.ObjectId_ = bson.NewObjectId()
	tn := utils.NowSecond()
	note.CreateTime = tn
	note.UpdateTime = tn

	if err := GetNoteCollection(db).Insert(note); err != nil {
		return err
	}

	return nil
}

//检查note是否存在
func (note *Note) IsExist() bool {
	noteTmp := NewNote(note.Path)
	noteTmp.UID = note.UID
	noteTmp, err := noteTmp.Get()

	if err == nil && noteTmp != nil {
		return true
	}

	return false
}

func (note *Note) updateBson(bsm bson.M) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()
	return GetNoteCollection(db).Update(bson.M{"UID": note.UID, "Path": note.Path}, bsm)
}

//更新笔记部分信息
func (note *Note) Update(fields []string) error {
	if len(fields) == 0 {
		return errors.New("len of fields is 0")
	}

	var fieldsMap bson.M = make(bson.M, 3)

	for _, v := range fields {
		fieldsMap[v] = reflect.ValueOf(*note).FieldByName(v).Interface()
	}

	return note.updateBson(bson.M{"$set": fieldsMap})
}

//更新笔记所有信息
func (note *Note) UpdateAll() error {
	var (
		data []byte
		err  error
		bsm  bson.M
	)

	if data, err = bson.Marshal(note); err != nil {
		return err
	}

	if err = bson.Unmarshal(data, &bsm); err != nil {
		return err
	}

	return note.updateBson(bsm)
}

//设置笔记状态
func (note *Note) SetStatus(status uint8) (err error) {
	return note.updateBson(bson.M{"$set": bson.M{"Status": status, "UpdateTime": note.UpdateTime}})
}

//设置笔记分享状态
func (note *Note) SetShare(status uint8) (err error) {
	return note.updateBson(bson.M{"$set": bson.M{"IsShare": status}})
}

//设置笔记投票状态
func (note *Note) SetVote(status uint8) (err error) {
	return note.updateBson(bson.M{"$set": bson.M{"IsVote": status}})
}

//设置笔记投票状态
func (note *Note) SetPasswordShare(status uint8) (err error) {
	return note.updateBson(bson.M{"$set": bson.M{"IsPassword": status}})
}

//更改笔记状态不可用
func (note *Note) Delete() error {
	return note.SetStatus(STATUS_DELETED)
}

//恢复笔记状态可用
func (note *Note) Recover() error {
	return note.SetStatus(STATUS_NORMAL)
}

//回收站删除笔记
func (note *Note) TrashClear() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetNoteCollection(db).Update(bson.M{"UID": note.UID, "Path": note.Path}, bson.M{"$set": bson.M{"Status": STATUS_TRASH}})
}

//彻底删除笔记
func (note *Note) Remove() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	if err := GetNoteCollection(db).Remove(bson.M{"UID": note.UID, "Path": note.Path}); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil
		}

		return err
	}

	return nil
}

//清除数据库数据
func (note *Note) Clear(updateTime time.Time) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	ch, err := GetNoteCollection(db).RemoveAll(bson.M{"Status": STATUS_TRASH, "UpdateTime": bson.M{"$lt": updateTime}})
	if err != nil {
		beego.Debug(ch)
	}
	return err
}

//移动笔记到另一个分组下
func (note *Note) Move(destCatePath, destCateName string) error {
	return note.updateBson(bson.M{"$set": bson.M{"CategoryPath": destCatePath, "CategoryName": destCateName}})
}
