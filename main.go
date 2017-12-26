package main

import (
	"notepad-api/conf"
	_ "notepad-api/conf"
	"notepad-api/model"
	_ "notepad-api/routers"
	"notepad-api/utils"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/redis.v5"
	//	"github.com/astaxie/beego/cache"
	//	_ "github.com/astaxie/beego/cache/redis"

	//	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	//	_ "github.com/go-sql-driver/mysql"
)

func init() {

	mongodb := beego.AppConfig.String("MONGODB_HOST")
	beego.Debug(mongodb)
	model.MSessionInit(mongodb, "", "", 100)
}

func main() {
	utils.Cache = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.RedisAddr,
		PoolSize: conf.Conf.PoolSize,
	})
	err := utils.Cache.Ping().Err()
	if err != nil {
		panic(err)
	}
	beego.SetStaticPath("/swagger", "swagger")
	beego.BConfig.WebConfig.Session.SessionOn = true
	tk1 := toolbox.NewTask("tk1", "0 0 1 * * 0-6", ClearAll)
	toolbox.AddTask("tk1", tk1)
	toolbox.StartTask()
	defer toolbox.StopTask()
	beego.Run()
}

func ClearAll() error {
	beego.Debug("Clear.........")
	var err error
	tm := time.Now().UTC()
	beego.Debug(tm)
	tm1 := tm.AddDate(0, 0, -30)
	note := model.NewNote("")
	category := model.NewCategory("", "")
	err = note.Clear(tm1)
	if err != nil {
		beego.Error(err)
	}
	err = category.Clear(tm1)
	if err != nil {
		beego.Error(err)
	}
	return err
}
