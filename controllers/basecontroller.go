package controllers

import (
	"encoding/json"
	cloudtoken "notepad-api/comm/token"
	"strings"
	"time"
	//	"encoding/json"
	"notepad-api/conf"
	"notepad-api/model"
	"notepad-api/operatelog"
	"notepad-api/utils"
	"strconv"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

var whiteNams = []string{"/API/V1.0/Share/Show", "/API/V1.0/Cloud/CallBack", "/API/V1.0/Cloud/GetDownloadURL", "/API/V1.0/Share/GetQuestion", "/API/V1.0/User/GetToken", "/API/V1.0/User/AppGetToken", "/API/V1.0/Vote/VoteGet", "/API/V1.0/Vote/ClickVote"}

func (this *BaseController) Prepare() {
	types := "web"
	AccountID := this.Ctx.Input.Header("AccountID")
	token := this.Ctx.Input.Header("token")
	version := this.Ctx.Input.Header("AppVer")
	//	Sign := this.Ctx.Input.Header("Sign")
	//	appID := this.Ctx.Input.Header("AppID")
	if len(version) != 0 {
		types = "app"
	}
	beego.Debug(AccountID, token)
	if conf.Conf.IsOpenLog {
		if this.Ctx.Request.RequestURI == "/API/V1.0/Share/Upload" {
			operatelog.SendToChannel(AccountID, this.Ctx.Request.RequestURI, types, "")
		} else {
			beego.Debug(string(this.Ctx.Input.RequestBody))
			operatelog.SendToChannel(AccountID, this.Ctx.Request.RequestURI, types, string(this.Ctx.Input.RequestBody))
		}
	}
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	this.Ctx.Output.Header("Access-Control-Allow-Headers", "AccountID,Content-Type,token")

	url := this.Ctx.Request.RequestURI
	if this.Ctx.Request.Method == "OPTIONS" {
		this.Ctx.WriteString("Success")
		return
	}
	versiontag, _ := strconv.Atoi(version)
	if versiontag < conf.Conf.Version {
		this.Data["json"] = utils.NOTE_LEVEV_NOT_SUPPORT
		this.ServeJSON()
		return
	}
	//	appkey := "40cb504c6e1e47b1b756caa541e9c522"
	//	if !utils.CheckAuthCodeSign(this.Ctx, string(this.Ctx.Input.RequestBody), appkey) {
	//		this.Data["json"] = utils.SIGN_CHECK_ERROR
	//		this.ServeJSON()
	//		return
	//	}
	beego.Debug(url)
	if !utils.CheckString(whiteNams, url) {
		if !utils.CheckToken(AccountID, token) && !cloudtoken.CheckToken(AccountID, token, utils.Cache) {
			beego.Debug("token失效")
			this.Data["json"] = utils.TOKEN_INVALID
			this.ServeJSON()
			return
		}
	}

}

func (this *BaseController) Finish() {
	url := this.Ctx.Request.RequestURI
	pushId := this.Ctx.Input.Header("PushID")
	//	if len(pushId) == 0 {
	//		return
	//	}
	if strings.ContainsAny(url, "Sync") {
		return
	}
	var registerID []string
	var syncMsg = make(map[string]interface{}, 0)
	syncMsg["MsgType"] = "startsync"
	syncMsg["Value"] = "NOTEPAD"
	syncMsg["Timestamp"] = time.Now().Unix()
	AccountID := this.Ctx.Input.Header("AccountID")
	syncMsg["Uid"] = AccountID
	resp := this.Data["json"]
	var temp map[string]interface{}
	var code int
	if bytes, err := json.Marshal(resp); err != nil {
		beego.Error(err)
	} else {
		json.Unmarshal(bytes, &temp)
	}
	if _, ok := temp["Code"].(float64); ok {
		code = int(temp["Code"].(float64))
	}
	//	appID := this.Ctx.Input.Header("AppID")

	if code == 0 {
		if data, ok := temp["Data"].(map[string]interface{}); ok {
			_, isHas := data["Path"]
			if _, ok := data["Usn"]; ok && isHas {
				syncMsg["Usn"] = data["Usn"]
				user := model.NewUserCollectionList(AccountID)
				if usercoll, err := user.GetList(); err == nil && len(usercoll) != 0 {
					for i := 0; i < len(usercoll); i++ {
						if !strings.EqualFold(pushId, usercoll[i].RegisterId) {
							registerID = append(registerID, usercoll[i].RegisterId)
						}
					}
				} else {
					return
				}
				if len(registerID) == 0 {
					return
				}
				msg, _ := json.Marshal(syncMsg)
				utils.PushCommand(registerID, string(msg), conf.Conf.PushId, conf.Conf.PushKey)
			}
		}
	}

}
