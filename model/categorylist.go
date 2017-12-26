package model

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type CategoryList struct {
	UID string
}

func NewCategoryList(UID string) *CategoryList {
	cl := new(CategoryList)
	cl.UID = UID

	return cl
}

func (cl *CategoryList) getListByBson(bsm bson.M) ([]Category, error) {

	var cateList []Category
	db := MSessionGet()
	if db == nil {
		return cateList, SessionFailError
	}
	defer db.Close()

	if err := GetCategoryCollection(db).Find(bsm).Sort("-CreateTime").All(&cateList); err != nil {

		return cateList, err
	}

	return cateList, nil
}

//获取用户下所有需要同步分组
func (cl *CategoryList) GetSynchList(afterUsn int) ([]Category, error) {
	var cateList []Category
	var bsm bson.M
	db := MSessionGet()
	if db == nil {
		return cateList, SessionFailError
	}
	defer db.Close()
	//	if afterUsn == 0 {
	//		bsm = bson.M{"Status": STATUS_NORMAL, "UID": cl.UID, "Usn": bson.M{"$gt": afterUsn}}
	//	} else {
	bsm = bson.M{"UID": cl.UID, "Usn": bson.M{"$gt": afterUsn}}
	//	}
	return cl.getListByBson(bsm)
}

//获取用户下所有正常分组
func (cl *CategoryList) GetList() ([]Category, error) {

	var cateList []Category
	db := MSessionGet()
	if db == nil {
		return cateList, SessionFailError
	}
	defer db.Close()

	return cl.getListByBson(bson.M{"UID": cl.UID, "Status": STATUS_NORMAL})
}

//模糊查询
func (nl *CategoryList) SearchCategoryList(title string) ([]Category, error) {
	var cateList []Category
	db := MSessionGet()
	if db == nil {
		return cateList, SessionFailError
	}
	defer db.Close()
	if err := GetCategoryCollection(db).Find(bson.M{"Name": bson.M{"$regex": title, "$options": "$i"}, "UID": nl.UID, "Status": STATUS_NORMAL}).All(&cateList); err != nil {
		return cateList, err
	}

	return cateList, nil
}

//获取用户下所有分组
func (cl *CategoryList) GetListAll() ([]Category, error) {
	var cateList []Category
	db := MSessionGet()
	if db == nil {
		return cateList, SessionFailError
	}
	defer db.Close()

	return cl.getListByBson(bson.M{"UID": cl.UID})
}

//获取用户下所有分组
func (cl *CategoryList) TrashAll() error {
	db := MSessionGet()
	if db == nil {
		return SessionFailError
	}
	defer db.Close()

	changeinfo, err := GetCategoryCollection(db).UpdateAll(bson.M{"UID": cl.UID, "Status": STATUS_DELETED}, bson.M{"$set": bson.M{"Status": STATUS_TRASH}})
	if err != nil {
		beego.Debug(changeinfo)
	}
	return err
}

//获取用户下所有已删除分组
func (cl *CategoryList) GetListDeleted() ([]Category, error) {
	var cateList []Category
	db := MSessionGet()
	if db == nil {
		return cateList, SessionFailError
	}
	defer db.Close()

	return cl.getListByBson(bson.M{"UID": cl.UID, "Status": STATUS_DELETED})

}

//获取用户下所有已删除分组
func (cl *CategoryList) GetListCount(status uint8) (int, error) {
	var (
		num int
		err error
	)
	db := MSessionGet()
	if db == nil {
		return 0, SessionFailError
	}
	defer db.Close()

	if num, err = GetCategoryCollection(db).Find(bson.M{"UID": cl.UID, "Status": status}).Count(); err == nil {
		return num, nil
	}
	return 0, err
}
