package operatelog

import (
	"time"

	"notepad-api/comm/mgodb"
	"notepad-api/conf"

	"github.com/astaxie/beego"
)

type OperateLog struct {
	UserId    string      `bson:"userId"`
	Router    string      `bson:"router"`
	Data      interface{} `bson:"data"`
	Type      string      `bson:"type"`
	TimeStamp int         `bson:"timeStamp"`
}

var channelNum, _ = beego.AppConfig.Int("channelNum")
var Sendc = make(chan map[string]interface{}, channelNum)

func init() {
	if conf.Conf.IsOpenLog {
		go SendToMgoDB()
	}
}

func SendToMgoDB() {
	for msg := range Sendc {
		beego.Debug("Send msg:", msg)
		var operatelog OperateLog
		operatelog.UserId = msg["userId"].(string)
		operatelog.Router = msg["router"].(string)
		operatelog.Data = msg["data"].(string)
		operatelog.TimeStamp = msg["time"].(int)
		operatelog.Type = msg["type"].(string)
		result := mgodb.AddRecord("cloudnote", "operate_log", operatelog)
		beego.Debug(result)
	}

}

//SendToChanner send message to channel
func SendToChannel(userId, router, types string, datas interface{}) {
	tm := int(time.Now().Unix())
	data := map[string]interface{}{
		"userId": userId,
		"router": router,
		"data":   datas,
		"time":   tm,
		"type":   types,
	}
	Sendc <- data
}
