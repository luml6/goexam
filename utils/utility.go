package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"notepad-api/conf"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/redis.v5"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"

	"github.com/satori/go.uuid"
	"qiniupkg.com/api.v7/kodo"
)

var (
	CATEGORY_NOT_EXIST  = &Response{Code: 1001, Message: "该笔记分组不存在", Data: nil}
	WITHOUT_PARAMETERS  = &Response{Code: 1002, Message: "参数不全", Data: nil}
	NOTE_NOT_EXIST      = &Response{Code: 1003, Message: "笔记不存在", Data: nil}
	NOTE_ALREADY_EXIST  = &Response{Code: 1004, Message: "该笔记或者附件已经存在，不能重复创建", Data: nil}
	CATEGORY_NAME_EXIST = &Response{Code: 1007, Message: "该笔记本名称已存在", Data: nil}
	//	NOTETITLE_ALREADY_EXIST   = &Response{Code: 1008, Message: "该笔记标题已存在", Data: nil}
	ATTACH_NOT_EXIST          = &Response{Code: 1009, Message: "附件不存在", Data: nil}
	ATTACH_UUID_ALREADY_EXIST = &Response{Code: 1010, Message: "附件Id已存在", Data: nil}
	NO_RIGHT_TO_ACCESS        = &Response{Code: 1011, Message: "没有权限访问", Data: nil}
	SHARE_ID_NOT_EXIST        = &Response{Code: 1012, Message: "分享ID不存在", Data: nil}
	SHARE_ANSWER_NOT_TRUE     = &Response{Code: 1013, Message: "分享答案不正确", Data: nil}
	NOTE_NOT_SHARE            = &Response{Code: 1014, Message: "该笔记没有被分享", Data: nil}
	NOTE_ALREADY_SHARE        = &Response{Code: 1015, Message: "该笔记已经被分享了", Data: nil}
	CATEGORY_NAME_NOT_NULL    = &Response{Code: 1016, Message: "该笔记分组不能为空", Data: nil}
	TOKEN_INVALID             = &Response{Code: 1017, Message: "业务token失效", Data: nil}
	CPTOKEN_INVALID           = &Response{Code: 1026, Message: "酷派token失效", Data: nil}
	USER_NOT_EXIST            = &Response{Code: 1018, Message: "用户不存在", Data: nil}
	PAPER_NOT_EXIST           = &Response{Code: 1019, Message: "该背景ID不存在", Data: nil}
	TEMPLATE_NOT_EXIST        = &Response{Code: 1020, Message: "模板ID不存在", Data: nil}
	WIDGET_NOT_EXIST          = &Response{Code: 1021, Message: "控件ID不存在", Data: nil}
	VOTE_ID_NOT_EXIST         = &Response{Code: 1022, Message: "投票ID不存在", Data: nil}
	NOTE_ALREADY_VOTE         = &Response{Code: 1023, Message: "该笔记已存在投票", Data: nil}
	OPTION_ID_NOT_EXIST       = &Response{Code: 1024, Message: "投票选项ID不存在", Data: nil}
	NOTE_NOT_VOTE             = &Response{Code: 1025, Message: "该笔记不存在投票", Data: nil}
	SHARE_ANSWER_NOT_NULL     = &Response{Code: 1027, Message: "分享答案不能为空", Data: nil}

	VOTE_QUESTION_NOT_NULL  = &Response{Code: 1028, Message: "投票问题不能为空", Data: nil}
	VOTE_OPTION_LIMIT       = &Response{Code: 1029, Message: "投票选项数量不能少于2个，大于20个", Data: nil}
	VOTE_OPTION_NOT_NULL    = &Response{Code: 1030, Message: "投票选项不能为空", Data: nil}
	VOTE_OPTION_NOT_SAME    = &Response{Code: 1031, Message: "投票选项不能相同", Data: nil}
	VOTE_OPTION_NOT_ONEMORE = &Response{Code: 1032, Message: "投票不支持多选", Data: nil}
	VOTE_OPTION_TIMEOUT     = &Response{Code: 1033, Message: "投票时间已过", Data: nil}
	SHARE_TYPE_NOT_EXIST    = &Response{Code: 1034, Message: "分享类型错误", Data: nil}
	SHARE_PASSWORD_NOT_NULL = &Response{Code: 1035, Message: "分享密码不能为空", Data: nil}
	SHARE_PASSWORD_LIMIT    = &Response{Code: 1036, Message: "分享答案为4~12位", Data: nil}
	SHARE_QUESTION_LIMIT    = &Response{Code: 1043, Message: "分享问题30位答案20位限制", Data: nil}
	SHARE_WRONG             = &Response{Code: 1039, Message: "分享失败", Data: nil}
	VOTE_TYPE_NOT_EXIST     = &Response{Code: 1037, Message: "投票类型错误", Data: nil}
	NOTE_TITLE_LIMIT        = &Response{Code: 1038, Message: "笔记title限制", Data: nil}
	NOTE_IS_LOCK            = &Response{Code: 1040, Message: "笔记被锁", Data: nil}
	NOTE_LEVEV_NOT_SUPPORT  = &Response{Code: 1041, Message: "笔记不支持版本过低", Data: nil}
	SIGN_CHECK_ERROR        = &Response{Code: 1042, Message: "sign校验失败", Data: nil}
	NOTE_NOT_CHANGE         = &Response{Code: 1044, Message: "该笔记没有更改", Data: nil}
	//PARAMSFORMATERROR 参数格式错误
	PARAMSFORMATERROR = &Response{Code: 1102, Message: "PARAMS_FORMAT_ERROR", Data: nil}
)
var Cache *redis.Client

const (
	SUCCESS_CODE = 0
	SYSTEM_CODE  = 1000
	SHAREKEY     = "NOTE:SHARE:"
	TOKENKEY     = "NOTE:TOKEN:"
	USNKEY       = "NOTE:USN:"
	NOTELOCKKEY  = "NOTE:LOCK:"
)

const ATTACH_PUBLIC_PREFIX string = "public:"

func HttpPost(url, params string, header http.Header) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(params))

	if err != nil {
		// handle error
		fmt.Println((err))
		return nil, err
	}
	req.Header = header

	resp, err := client.Do(req)
	if err != nil {
		// handle error
		fmt.Println((err))
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)

}

//NewResponse 初始化response
func NewResponse(Code int, Msg string, data interface{}) *Response {
	if Code != 0 {
		data = nil
	}
	return &Response{
		Code:    Code,
		Message: Msg,
		Data:    data,
	}
}

//CreateId create one id
func CreateId() string {
	return uuid.NewV4().String()
}

//CreateId create one id
func CreateRandId(leng int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < leng; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func CheckString(list []string, st string) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == st {
			return true
		}
	}
	return false
}

func GetUserInfo(uid, accesstoken, refreshtoken, appid, cpurl string) *UserInfo {
	var userinfo UserInfo
	u, _ := url.Parse(cpurl)
	q := u.Query()
	q.Set("access_token", accesstoken)
	q.Set("oauth_consumer_key", appid)
	q.Set("openid", uid)
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		beego.Error(err)
		userinfo.Rtncode = "-1"
		return &userinfo
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		beego.Error(err)
		userinfo.Rtncode = "-1"
		return &userinfo
	}
	beego.Debug(string(result))
	json.Unmarshal(result, &userinfo)
	return &userinfo
}

func AppCheck(uid, tkt, appid, cpurl string) *AppResp {
	var userinfo AppResp
	resp, err := http.PostForm(cpurl,
		url.Values{"tkt": {tkt}, "appid": {appid}, "uid": {uid}})

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		beego.Error(err)
		userinfo.Rtncode = "-1"
		return &userinfo
	}
	beego.Debug(string(body))
	json.Unmarshal(body, &userinfo)
	return &userinfo
}

func RefreshToken(appkey, refreshtoken, appid, cpurl string) *RefreshResp {
	var resp RefreshResp
	u, _ := url.Parse(cpurl)
	q := u.Query()
	q.Set("grant_type", "refresh_token")
	q.Set("client_id", appid)
	q.Set("refresh_token", refreshtoken)
	q.Set("client_secret", appkey)
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		beego.Error(err)
		return &resp
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		beego.Error(err)
		return &resp
	}
	json.Unmarshal(result, &resp)
	return &resp
}
func CreateToken(uid string, tm int) (error, string) {
	h := md5.New()
	timestamp := time.Now().Unix()
	code := fmt.Sprintf("%v-%v", uid, timestamp)
	h.Write([]byte(code))
	token := hex.EncodeToString(h.Sum(nil))
	key := TOKENKEY + uid
	keepTime := time.Duration(tm) * time.Hour
	//存入redis
	err := Cache.Set(key, token, keepTime).Err()
	beego.Debug(Cache.Get(key).Result())
	//判断token是否失效
	if err != nil {
		return err, ""
	}
	return nil, token
}

func ShareCache(key, value string, tm int) error {
	keepTime := time.Duration(tm) * time.Hour
	//存入redis
	err := Cache.Set(key, value, keepTime).Err()

	return err
}

func LockCache(key string, value interface{}, tm int) error {

	keepTime := time.Duration(tm) * time.Second
	//存入redis
	err := Cache.Set(key, value, keepTime).Err()

	return err
}

func CheckStingsNotNil(list []interface{}) bool {
	for i := 0; i < len(list); i++ {
		str, _ := list[i].(string)
		if str == "" {
			return true
		}
	}
	return false
}

func CheckSlice(list *[]interface{}) bool {
	var x []interface{} = []interface{}{}
	for _, i := range *list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					return true
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return false
}

func CheckToken(uid, token string) bool {
	key := TOKENKEY + uid
	if user, err := Cache.Get(key).Result(); err == nil {
		if user == token {
			return true
		}
	}
	if user, err := Cache.Get(token).Result(); err == nil {
		if user == uid {
			return true
		}
	} else {
		beego.Error(err)
	}

	return false
}

func CheckLock(key string) bool {
	if ok, err := Cache.Get(key).Result(); len(ok) != 0 && err == nil {
		return true
	}
	return false
}

func NowSecond() time.Time {
	return time.Unix(0, time.Now().UTC().UnixNano()/1e9*1e9)
}

func GetQiniuhash(bucket, uuid string) string {
	var qiniu map[string]interface{}
	var hash string
	c := kodo.New(0, nil)
	encodedEntryURI := conf.Conf.MessageUrl + kodo.URIStat(bucket, uuid)
	resp, err := c.Get(encodedEntryURI)
	if err != nil {
		beego.Error(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &qiniu)
	if _, ok := qiniu["hash"].(string); ok {
		hash = qiniu["hash"].(string)
	}
	return hash
}

func CheckAuthCodeSign(ctx *context.Context, body string, appkey string) bool {
	signMap := make(map[string]interface{})
	signMap["SerVer"] = ctx.Input.Header("SerVer")
	signMap["AppID"] = ctx.Input.Header("AppID")
	signMap["AppVer"] = ctx.Input.Header("AppVer")
	signMap["Body"] = body
	fmt.Println(signMap)
	sign, _ := GenSign(signMap, appkey)
	headsign := ctx.Input.Header("Sign")
	if headsign != sign {
		beego.Debug(headsign, sign)
		beego.Error("Check Sign error,head sign is :%s,calc sign is :%s", headsign, sign)
		return false //for test return true
	}
	return true

}

func GetTime(tm interface{}) time.Time {
	if update, ok := tm.(float64); ok && update != 0 {
		return time.Unix(int64(update), 0)
	}
	return NowSecond()
}

func PushCommand(alias []string, msg string, appid, appkey string) error {
	header := make(http.Header)
	ts := time.Now().Unix()
	timestamp := strconv.FormatInt(ts, 10)

	var payload pushPayload
	payload.Environment = "production"
	payload.ApiType = "alias"
	payload.TimeType = 0
	payload.Expires = ts + 3600
	payload.IsCallback = false
	payload.Body = msg
	payload.PushType = "pass_through"
	payload.Target = alias

	var body string

	if json_byte, err := json.Marshal(payload); err != nil {
		beego.Error(err)
		return errors.New("GetUDID json marshal error")
	} else {
		body = string(json_byte)
	}

	beego.Debug(body)

	header.Set("Connection", "keep-alive")
	header.Set("Timestamp", timestamp)
	header.Set("AppID", appid)
	beego.Debug("header:%s", header)
	header.Set("Sign", Sign(appid, appkey, timestamp, body))

	if body_byte, err := HttpPost(conf.Conf.PushUrl, body, header); err != nil {
		return err
	} else {
		beego.Debug(string(body_byte))
	}

	return nil
}

/*计算sign签名*/
func Sign(appid, appkey, timestamp, body string) string {
	dict := make(map[string]string)
	dict["AppID"] = appid
	dict["Timestamp"] = timestamp
	dict["Body"] = body
	dict["AppKey"] = appkey

	keys := make([]string, len(dict))
	i := 0
	for k, _ := range dict {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	signStr := ""

	for _, k := range keys {
		signStr += k + dict[k]
	}
	beego.Debug(signStr)

	t := sha1.New()
	io.WriteString(t, signStr)
	sign := hex.EncodeToString(t.Sum(nil))

	return sign
}
