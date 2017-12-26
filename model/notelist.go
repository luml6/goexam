package model

import (
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type NoteList struct {
	CategoryPath string
	UID          string
}

func NewNoteList(catePath, uid string) *NoteList {
	nl := new(NoteList)
	nl.CategoryPath = catePath
	nl.UID = uid

	return nl
}

func (nl *NoteList) getListByBson(bsm bson.M) ([]Note, error) {
	var noteList []Note
	db := MSessionGet()
	if db == nil {
		return noteList, SessionFailError
	}
	defer db.Close()

	bsm["UID"] = nl.UID

	if err := GetNoteCollection(db).Find(bsm).Sort("-IsPinned", "-UpdateTime").All(&noteList); err != nil {
		return noteList, err
	}

	return noteList, nil
}

//获取用户所有正常笔记列表
func (nl *NoteList) GetAllList() ([]Note, error) {
	var noteList []Note
	db := MSessionGet()
	if db == nil {
		return noteList, SessionFailError
	}
	defer db.Close()

	if err := GetNoteCollection(db).Find(bson.M{"UID": nl.UID, "Status": STATUS_NORMAL}).Sort("-IsPinned", "-UpdateTime").All(&noteList); err != nil {
		return noteList, err
	}

	return noteList, nil
}

//模糊查询
func (nl *NoteList) SearchNoteList(title string) ([]Note, error) {
	var noteList []Note
	db := MSessionGet()
	if db == nil {
		return noteList, SessionFailError
	}
	defer db.Close()
	if err := GetNoteCollection(db).Find(bson.M{"Title": bson.M{"$regex": title, "$options": "$i"}, "UID": nl.UID, "Status": STATUS_NORMAL}).Sort("-IsPinned", "-UpdateTime").All(&noteList); err != nil {
		return noteList, err
	}

	return noteList, nil
}

//模糊查询
func (nl *NoteList) SearchNoteByCatePath(title, categoryPath string) ([]Note, error) {
	var noteList []Note
	db := MSessionGet()
	if db == nil {
		return noteList, SessionFailError
	}
	defer db.Close()
	if err := GetNoteCollection(db).Find(bson.M{"UID": nl.UID, "CategoryPath": categoryPath, "Title": bson.M{"$regex": title, "$options": "$i"}, "Status": STATUS_NORMAL}).Sort("-IsPinned", "-UpdateTime").All(&noteList); err != nil {
		return noteList, err
	}

	return noteList, nil
}

//获取用户所有笔记列表
func (nl *NoteList) GetNoteListCount(status uint8) (int, error) {
	var num int
	var err error
	db := MSessionGet()
	if db == nil {
		return 0, SessionFailError
	}
	defer db.Close()
	if num, err = GetNoteCollection(db).Find(bson.M{"UID": nl.UID, "Status": status}).Count(); err != nil {
		return 0, err
	}
	return num, err
}

//获取用户所有笔记列表
func (nl *NoteList) GetNoteList(bsm bson.M) ([]Note, error) {
	var noteList []Note
	db := MSessionGet()
	if db == nil {
		return noteList, SessionFailError
	}
	defer db.Close()
	bsm["UID"] = nl.UID
	if err := GetNoteCollection(db).Find(bsm).Sort("-IsPinned", "-UpdateTime").All(&noteList); err != nil {
		return noteList, err
	}

	return noteList, nil
}

//获取用户所有笔记列表
func (nl *NoteList) GetSyncNoteList(bsm bson.M) ([]Note, error) {
	var noteList []Note
	db := MSessionGet()
	if db == nil {
		return noteList, SessionFailError
	}
	defer db.Close()
	bsm["UID"] = nl.UID
	if err := GetNoteCollection(db).Find(bsm).Select(bson.M{"Path": 1, "Status": 1, "CheckSum": 1, "CategoryPath": 1, "Title": 1, "Usn": 1, "UpdateTime": 1, "IsPinned": 1, "IsVote": 1, "IsPassworded": 1}).All(&noteList); err != nil {
		return noteList, err
	}

	return noteList, nil
}

//获取用户所有删除笔记列表
func (nl *NoteList) GetNoteDeletedList(bsm bson.M) ([]Note, error) {
	var noteList []Note
	db := MSessionGet()
	if db == nil {
		return noteList, SessionFailError
	}
	defer db.Close()
	bsm["UID"] = nl.UID
	bsm["Status"] = STATUS_DELETED
	if err := GetNoteCollection(db).Find(bsm).Sort("-UpdateTime").All(&noteList); err != nil {
		return noteList, err
	}

	return noteList, nil
}

//获取正常笔记列表
func (nl *NoteList) GetList() ([]Note, error) {
	return nl.getListByBson(bson.M{"UID": nl.UID, "CategoryPath": nl.CategoryPath, "Status": STATUS_NORMAL})
}

//获取删除笔记列表
func (nl *NoteList) GetDeleted() ([]Note, error) {
	return nl.getListByBson(bson.M{"UID": nl.UID, "CategoryPath": nl.CategoryPath, "Status": STATUS_DELETED})
}

//获取笔记列表
func (nl *NoteList) GetListNoCategory() ([]Note, error) {
	return nl.getListByBson(bson.M{"UID": nl.UID, "CategoryPath": "", "Status": STATUS_NORMAL})
}

//删除category属性
func (nl *NoteList) DeleteCategory() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	if _, err := GetNoteCollection(db).UpdateAll(bson.M{"UID": nl.UID, "CategoryPath": nl.CategoryPath},
		bson.M{"$set": bson.M{"CategoryPath": "", "CategoryName": ""}}); err != nil {

		return err
	}

	return nil
}

//更新categoryName属性
func (nl *NoteList) UpdateCategory(name string) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	if _, err := GetNoteCollection(db).UpdateAll(bson.M{"UID": nl.UID, "CategoryPath": nl.CategoryPath},
		bson.M{"$set": bson.M{"CategoryName": name}}); err != nil {

		return err
	}

	return nil
}

//删除note
func (nl *NoteList) Delete(updateTime time.Time) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	if _, err := GetNoteCollection(db).UpdateAll(bson.M{"CategoryPath": nl.CategoryPath},
		bson.M{"$set": bson.M{"Status": STATUS_DELETED, "UpdateTime": updateTime}}); err != nil {
		return err
	}

	return nil
}

//恢复笔记
func (nl *NoteList) Recover() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	if _, err := GetNoteCollection(db).UpdateAll(bson.M{"CategoryPath": nl.CategoryPath},
		bson.M{"$set": bson.M{"Status": STATUS_NORMAL}}); err != nil {
		return err
	}

	return nil
}

//获取用户下所有分组
func (nl *NoteList) TrashAll() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	changeinfo, err := GetNoteCollection(db).UpdateAll(bson.M{"UID": nl.UID, "Status": STATUS_DELETED}, bson.M{"$set": bson.M{"Status": STATUS_TRASH}})
	if err != nil {
		beego.Debug(changeinfo)
	}
	return err
}

//彻底删除笔记
func (nl *NoteList) TrashClear() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	if _, err := GetNoteCollection(db).UpdateAll(bson.M{"CategoryPath": nl.CategoryPath, "Status": STATUS_DELETED},
		bson.M{"$set": bson.M{"Status": STATUS_TRASH}}); err != nil {
		return err
	}

	return nil
}
