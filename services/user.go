package services

import (
	cloudtoken "notepad-api/comm/token"
	"notepad-api/conf"
	"notepad-api/model"
	"notepad-api/utils"
	"strings"
	"time"

	"github.com/astaxie/beego"
	//	"gopkg.in/mgo.v2/bson"
)

// 自增Usn
// 每次notebook,note添加, 修改, 删除, 都要修改
func IncrUsn(userId string) int {
	var usn int
	filed := []string{"Usn", "UpdateTime"}
	user := model.NewUser(userId)
	if _, err := user.Get(); err == nil {
		usn = user.Usn
	} else {
		beego.Error(err)
		return 0
	}
	usn += 1
	user.Usn = usn
	user.UpdateTime = time.Now().UTC()
	beego.Debug("inc Usn")
	if err := user.Update(filed); err != nil {
		beego.Error(err)
	}
	key := utils.USNKEY + userId
	if err := utils.Cache.Set(key, usn, 0).Err(); err != nil {
		beego.Error(err)
	}
	return usn
}

func AppGetToken(Ob map[string]interface{}) *utils.Response {
	var (
		Uid        string
		registerId string
		AuthToken  string
		token      string
		//		userInfo   *utils.AppResp
		err    error
		islock bool
	)

	if _, ok := Ob["Uid"].(string); ok {
		Uid = Ob["Uid"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["AuthToken"].(string); ok {
		AuthToken = Ob["AuthToken"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["RegisterID"].(string); ok {
		registerId = Ob["RegisterID"].(string)
	}
	//	userInfo = utils.AppCheck(Uid, AuthToken, conf.Conf.AppID, conf.Conf.AppUrl)
	beego.Debug(AuthToken)
	key := utils.TOKENKEY + Uid
	//	if userInfo.Rtncode == "0" {
	if token, err = utils.Cache.Get(key).Result(); len(token) == 0 || err != nil {
		//			token = string(value)
		err, token = utils.CreateToken(Uid, conf.Conf.CacheTime)
	}
	if err != nil || token == "" {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	user := model.NewUser(Uid)
	userTmp, err := user.Get()
	if userTmp == nil || err != nil {
		if err := user.Add(); err != nil {
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}
	}
	//		} else {
	//			return utils.CPTOKEN_INVALID
	//		}
	t := struct {
		Uid         string
		AccessToken string
	}{Uid, token}
	if len(registerId) != 0 {
		Uidkey := utils.NOTELOCKKEY + Uid
		if utils.CheckLock(Uidkey) {
			return utils.NewResponse(0, "", t)
		}
		utils.LockCache(Uidkey, 1, 30)
		islock = true
		register := model.NewUserCollection(Uid)
		register.RegisterId = registerId
		register.Status = model.STATUS_NORMAL
		tn := utils.NowSecond()
		register.UpdateTime = tn
		if users, err := register.GetByRegisterID(); err == nil && users != nil {
			filed := []string{"Status", "RegisterId", "UpdateTime"}
			if !strings.EqualFold(users.UID, Uid) {
				if err := register.Update(filed); err != nil {
					beego.Error(err)
				}
			}
		} else {
			beego.Error(err)
			if err := register.Add(); err != nil {
				beego.Error(err)
			}
		}
		if islock {
			utils.Cache.Del(Uidkey)
		}
	}

	return utils.NewResponse(0, "", t)
}
func PushAccount(Ob map[string]interface{}) *utils.Response {
	var (
		Uid          string
		AccessToken  string
		RefreshToken string
		token        string
		authcode     string
		file         []string
		nikeName     string
		headImg      string
		//		userInfo     *utils.UserInfo
		err error
	)
	if _, ok := Ob["Uid"].(string); ok {
		Uid = Ob["Uid"].(string)
	}
	if _, ok := Ob["AccessToken"].(string); ok {
		AccessToken = Ob["AccessToken"].(string)
		file = append(file, "AccessToken")
	}
	if _, ok := Ob["RefreshToken"].(string); ok {
		RefreshToken = Ob["RefreshToken"].(string)
		file = append(file, "RefreshToken")
	}
	if authcode, token, err = cloudtoken.GetCloudToken(Uid, AccessToken, conf.Conf.AppID, conf.Conf.UserUrl, conf.Conf.CacheTime, utils.Cache); len(token) != 0 && err == nil {
		user := model.NewUser(Uid)
		userTmp, err := user.Get()
		if userTmp == nil || err != nil {
			user.AccessToken = AccessToken
			user.RefreshToken = RefreshToken
			if err := user.Add(); err != nil {
				return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
			}
		} else {
			userTmp.AccessToken = AccessToken
			userTmp.RefreshToken = RefreshToken
			if err := userTmp.Update(file); err != nil {
				return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
			}
		}
	} else {
		return utils.CPTOKEN_INVALID
	}
	userInfo := utils.GetUserInfo(Uid, AccessToken, RefreshToken, conf.Conf.AppID, conf.Conf.CpUserUrl)
	if userInfo.Rtncode == "0" {
		nikeName = userInfo.Nickname
		headImg = userInfo.HeadIconUrl
	}
	//		if token, err = utils.Cache.Get(key).Result(); len(token) == 0 || err != nil {
	//			err, token = utils.CreateToken(Uid, conf.Conf.CacheTime)
	//		}
	//		if err != nil || token == "" {
	//			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	//		}
	//		user := model.NewUser(Uid)
	//		userTmp, err := user.Get()
	//		if userTmp == nil || err != nil {
	//			user.AccessToken = AccessToken
	//			user.RefreshToken = RefreshToken
	//			if err := user.Add(); err != nil {
	//				return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	//			}
	//		} else {
	//			userTmp.AccessToken = AccessToken
	//			userTmp.RefreshToken = RefreshToken
	//			if err := userTmp.Update(file); err != nil {
	//				return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	//			}
	//		}
	//	} else {
	//		return utils.CPTOKEN_INVALID
	//	}
	t := struct {
		AuthCode string
		Uid      string
		Token    string
		NickName string
		HeadImg  string
	}{authcode, Uid, token, nikeName, headImg}
	return utils.NewResponse(0, "", t)
}

func GetUserInfo(userId string) *utils.Response {
	var file = []string{"AccessToken", "RefreshToken", "UpdateTime"}
	user := model.NewUser(userId)
	userTmp, err := user.Get()
	if userTmp != nil && err == nil {
		beego.Debug(userTmp)
		userInfo := utils.GetUserInfo(userId, userTmp.AccessToken, userTmp.RefreshToken, conf.Conf.AppID, conf.Conf.UserUrl)
		if userInfo.Rtncode == "0" {
			return utils.NewResponse(0, "", userInfo)
		} else if userInfo.Rtncode == "2014" && userTmp.RefreshToken != "" {
			refresh := utils.RefreshToken(conf.Conf.Appkey, userTmp.RefreshToken, conf.Conf.AppID, conf.Conf.RefreshUrl)
			userInfo = utils.GetUserInfo(userId, refresh.AccessToken, refresh.RefreshToken, conf.Conf.AppID, conf.Conf.UserUrl)
			if userInfo.Rtncode != "0" {
				return utils.CPTOKEN_INVALID
			}
			userTmp.AccessToken = refresh.AccessToken
			userTmp.RefreshToken = refresh.RefreshToken
			userTmp.UpdateTime = utils.NowSecond()
			if err := user.Update(file); err != nil {
				return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
			}
			return utils.NewResponse(0, "", userInfo)
		}
		return utils.CPTOKEN_INVALID
	}
	return utils.USER_NOT_EXIST
}
func AllTrash(userId string) *utils.Response {
	//	cates := model.NewCategoryList(userId)
	//	if err := cates.TrashAll(); err != nil {
	//		beego.Debug(err)
	//		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	//	}
	noteList := model.NewNoteList("", userId)
	if err := noteList.TrashAll(); err != nil {
		beego.Debug(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		IsSuccess bool
	}{true}
	return utils.NewResponse(0, "", t)
}
func LoginOut(userId, token string) bool {
	err := utils.Cache.Del(token).Err()
	if err != nil {
		return false
	}

	return true
}

func Cancellation(userId, registerId string) bool {
	filed := []string{"Status", "UpdateTime"}
	usercoll := model.NewUserCollection(userId)
	usercoll.RegisterId = registerId
	coll, err := usercoll.Get()
	if coll == nil && err != nil {
		beego.Error(err)
		return false
	}
	coll.Status = model.STATUS_DELETED
	tn := utils.NowSecond()
	coll.UpdateTime = tn
	if err := coll.Update(filed); err != nil {
		beego.Error(err)
		return false
	}
	return true
}
