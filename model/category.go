package model

import (
	"errors"
	"reflect"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//分组信息
type Category struct {
	ObjectId_  bson.ObjectId `bson:"_id"`
	Path       string        `bson:"Path"`       //分组唯一路径
	Name       string        `bson:"Name"`       //分组名称
	UID        string        `bson:"UID"`        //用户唯一ID
	Status     uint8         `bson:"Status"`     //分组状态
	CreateTime time.Time     `bson:"CreateTime"` //创建时间
	UpdateTime time.Time     `bson:"UpdateTime"` //更新时间
	Source     uint8         `bson:"Source"`     //笔记本来源
	//	IsTrash    uint8         `bson:"IsTrash"`    //是否回收站删除
	Usn int `bson:"Usn"` //笔记更新时间
}

func NewCategory(UID, path string) *Category {
	cate := new(Category)
	cate.UID = UID
	cate.Path = path

	return cate
}

//添加分组
func (cate *Category) Add() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	cate.ObjectId_ = bson.NewObjectId()
	//	tn := utils.NowSecond()

	//	cate.CreateTime = tn
	//	cate.UpdateTime = tn
	cate.Status = STATUS_NORMAL

	if err := GetCategoryCollection(db).Insert(cate); err != nil {
		return err
	}

	return nil
}

func (cate *Category) getBson(bsm bson.M) (*Category, error) {
	db := MSessionGet()
	if db == nil {
		return cate, SessionFailError
	}
	defer db.Close()

	if err := GetCategoryCollection(db).Find(bsm).One(cate); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil, nil
		}

		return cate, err
	}

	return cate, nil
}

//根据分组名称和用户id查询
func (cate *Category) GetByNameAndUID(name, UID string) (*Category, error) {
	return cate.getBson(bson.M{"UID": UID, "Name": name, "Status": STATUS_NORMAL})
}

//根据分组名称和用户id查询回收站
func (cate *Category) GetRecycleByNameAndUID(name, UID string) (*Category, error) {
	return cate.getBson(bson.M{"UID": UID, "Name": name, "Status": STATUS_DELETED})
}

//根据分组路径和状态查询
func (cate *Category) GetByPathAndStatus(path string, status uint8) (*Category, error) {
	return cate.getBson(bson.M{"UID": cate.UID, "Path": path, "Status": status})
}

//检查分组是否存在
func (cate *Category) IsExist() bool {
	cateTmp := NewCategory(cate.UID, cate.Path)
	cateTmp, err := cateTmp.Get()

	if err == nil && cateTmp != nil {
		return true
	}

	return false
}

//查询可用分组
func (cate *Category) Get() (*Category, error) {
	return cate.getBson(bson.M{"UID": cate.UID, "Path": cate.Path, "Status": STATUS_NORMAL})
}

//查询所有状态分组
func (cate *Category) GetAllStatus() (*Category, error) {
	return cate.getBson(bson.M{"UID": cate.UID, "Path": cate.Path})
}

func (cate *Category) updateBson(bsm bson.M) (err error) {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetCategoryCollection(db).Update(bson.M{"UID": cate.UID, "Path": cate.Path}, bsm)
}

//更新分组部分信息
func (cate *Category) Update(fields []string) (err error) {
	if len(fields) == 0 {
		return errors.New("len of fields is 0")
	}

	var fieldsMap bson.M = make(bson.M, 3)

	for _, v := range fields {
		fieldsMap[v] = reflect.ValueOf(*cate).FieldByName(v).Interface()
	}

	return cate.updateBson(bson.M{"$set": fieldsMap})
}

//更新分组所有信息
func (cate *Category) UpdateAll() error {
	var (
		data []byte
		err  error
		bsm  bson.M
	)

	if data, err = bson.Marshal(cate); err != nil {
		return err
	}

	if err = bson.Unmarshal(data, &bsm); err != nil {
		return err
	}

	return cate.updateBson(bsm)
}

//更寻分组状态
func (cate *Category) updateStatus(status uint8) (err error) {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()
	return GetCategoryCollection(db).Update(bson.M{"UID": cate.UID, "Path": cate.Path}, bson.M{"$set": bson.M{"Status": status, "UpdateTime": cate.UpdateTime}})
}

//更改分组状态不可用
func (cate *Category) Delete() (err error) {
	//TODO 分组所有笔记修改分组属性，改为未分组
	return cate.updateStatus(STATUS_DELETED)
}

//更改分组状态可用
func (cate *Category) Recover() error {
	return cate.updateStatus(STATUS_NORMAL)
}

//回收站删除分组
func (cate *Category) TrashClear() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetCategoryCollection(db).Update(bson.M{"Path": cate.Path}, bson.M{"$set": bson.M{"Status": STATUS_TRASH}})
}

//彻底删除分组
func (cate *Category) Remove() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	return GetCategoryCollection(db).Remove(bson.M{"Path": cate.Path})
}

//清除数据库数据
func (cate *Category) Clear(updateTime time.Time) error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()
	ch, err := GetCategoryCollection(db).RemoveAll(bson.M{"Status": STATUS_TRASH, "UpdateTime": bson.M{"$lt": updateTime}})
	if err != nil {
		beego.Debug(ch)
	}
	return err
}
