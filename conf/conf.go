package conf

import (
	"github.com/astaxie/beego"
)

type Config struct {
	UserUrl     string
	CpUserUrl   string
	AppID       string
	RefreshUrl  string
	Appkey      string
	AppUrl      string
	CacheTime   int
	ShareUrl    string
	IsOpenLog   bool
	LoginOutUrl string
	DownUrl     string
	MessageUrl  string
	PhotoStyle  string
	PoolSize    int
	RedisAddr   string
	Version     int
	PushUrl     string
	PushId      string
	PushKey     string
}

var Conf *Config

func init() {
	Conf = new(Config)
	Conf.DownUrl = beego.AppConfig.String("download")
	Conf.UserUrl = beego.AppConfig.String("userurl")
	Conf.CpUserUrl = beego.AppConfig.String("cp::GetUserUrl")
	Conf.AppID = beego.AppConfig.String("cp::AppID")
	Conf.RefreshUrl = beego.AppConfig.String("cp::RefreshUrl")
	Conf.Appkey = beego.AppConfig.String("cp::appkey")
	Conf.AppUrl = beego.AppConfig.String("cp::appurl")
	Conf.CacheTime, _ = beego.AppConfig.Int("CacheTime")
	Conf.ShareUrl = beego.AppConfig.String("sharUrl")
	Conf.LoginOutUrl = beego.AppConfig.String("loginouturl")
	Conf.IsOpenLog, _ = beego.AppConfig.Bool("IsOpenLog")
	Conf.MessageUrl = beego.AppConfig.String("qiniu::messageurl")
	Conf.PhotoStyle = beego.AppConfig.String("qiniu::style")
	Conf.PoolSize, _ = beego.AppConfig.Int("poolsize")
	Conf.RedisAddr = beego.AppConfig.String("redisAddr")
	Conf.Version, _ = beego.AppConfig.Int("version")
	Conf.PushUrl = beego.AppConfig.String("pushurl")
	Conf.PushId = beego.AppConfig.String("push::appid")
	Conf.PushKey = beego.AppConfig.String("push::appkey")
}
