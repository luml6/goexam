package mgodb

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var URL = beego.AppConfig.String("MONGODB_HOST") //mongodb连接字符串

var (
	mgoSession *mgo.Session
)

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(URL)
		if err != nil {
			panic(err) //直接终止程序运行
		}
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

//公共方法，获取collection对象
func WitchCollection(dataBase string, collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}

//用于两个返回值的函数: RemoveAll, UpdateAll
func WitchCollectionAll(dataBase string, collection string, s func(*mgo.Collection) (interface{}, error)) (interface{}, error) {
	session := getSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}

//公共方法，获取collection后的count值
func WitchCollectionCount(dataBase string, collection string, s func(*mgo.Collection) (int, error)) (int, error) {
	session := getSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}

func AddRecord(dataBase string, collection string, p interface{}) string {
	insert := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := WitchCollection(dataBase, collection, insert)
	if err != nil {
		beego.Debug(err)
		return "false"
	}
	return "true"
}

func RemoveData(dataBase string, collection string, query interface{}) string {
	removefunc := func(c *mgo.Collection) error {
		return c.Remove(query)
	}
	err := WitchCollection(dataBase, collection, removefunc)
	if err != nil {
		beego.Info("RemoveAllData Error is: ", err)
		return "false"
	}
	return "true"
}

func RemoveAllData(dataBase string, collection string, query interface{}) string {
	removefunc := func(c *mgo.Collection) (interface{}, error) {
		return c.RemoveAll(query)
	}
	_, err := WitchCollectionAll(dataBase, collection, removefunc)
	if err != nil {
		beego.Info("RemoveAllData Error is: ", err)
		return "false"
	}
	return "true"
}

func UpdateData(dataBase string, collection string, selector interface{}, update interface{}) string {
	updatefunc := func(c *mgo.Collection) error {
		return c.Update(selector, update)
	}
	err := WitchCollection(dataBase, collection, updatefunc)
	if err != nil {
		beego.Info("UpdateData Error is: ", err)
		return "false"
	}
	return "true"
}

func UpdateAllData(dataBase string, collection string, selector interface{}, update interface{}) string {
	updatefunc := func(c *mgo.Collection) (interface{}, error) {
		return c.UpdateAll(selector, update)
	}
	_, err := WitchCollectionAll(dataBase, collection, updatefunc)
	if err != nil {
		beego.Info("UpdateAllData Error is: ", err)
		return "false"
	}
	return "true"
}

/***/

func AddRecordNew(dataBase string, collection string, p interface{}) (result bool, err_msg string) {
	insert := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := WitchCollection(dataBase, collection, insert)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func RemoveDataNew(dataBase string, collection string, query interface{}) (result bool, err_msg string) {
	removefunc := func(c *mgo.Collection) error {
		return c.Remove(query)
	}
	err := WitchCollection(dataBase, collection, removefunc)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func RemoveAllDataNew(dataBase string, collection string, query interface{}) (result bool, err_msg string) {
	removefunc := func(c *mgo.Collection) (interface{}, error) {
		return c.RemoveAll(query)
	}
	_, err := WitchCollectionAll(dataBase, collection, removefunc)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func UpdateDataNew(dataBase string, collection string, selector interface{}, update interface{}) (result bool, err_msg string) {
	updatefunc := func(c *mgo.Collection) error {
		return c.Update(selector, update)
	}
	err := WitchCollection(dataBase, collection, updatefunc)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func UpdateAllDataNew(dataBase string, collection string, selector interface{}, update interface{}) (result bool, err_msg string) {
	updatefunc := func(c *mgo.Collection) (interface{}, error) {
		return c.UpdateAll(selector, update)
	}
	_, err := WitchCollectionAll(dataBase, collection, updatefunc)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

/*
//更新person数据
func UpdatePerson(query bson.M, change bson.M) string {
	exop := func(c *mgo.Collection) error {
		return c.Update(query, change)
	}
	err := witchCollection("person", exop)
	if err != nil {
		return "true"
	}
	return "false"
}
*/

/**
 * [Search description]
 * @param {[type]} database		  string [description]
 * @param {[type]} collectionName string [description]
 * @param {[type]} query          bson.M [description])   (count int, err error [description]
 */
func SearchCount(dataBase string, collectionName string, query bson.M) (int, error) {
	exop := func(c *mgo.Collection) (int, error) {
		return c.Find(query).Count()
	}
	return WitchCollectionCount(dataBase, collectionName, exop)
}

/**
 * [SearchRecords description]
 * @param {[type]} collectionName string [description]
 * @param {[type]} query          bson.M [description]
 * @param {[type]} sort           string [description]
 * @param {[type]} fields         bson.M [description]
 * @param {[type]} skip           int    [description]
 * @param {[type]} limit          int)   (results      []interface{}, err error [description]
 */
func SearchRecords(dataBase string, collectionName string, query bson.M, sort string, fields bson.M, skip int, limit int) (results []interface{}, err error) {
	exop := func(c *mgo.Collection) error {
		return c.Find(query).Sort(sort).Select(fields).Skip(skip).Limit(limit).All(&results)
	}
	err = WitchCollection(dataBase, collectionName, exop)
	return
}
