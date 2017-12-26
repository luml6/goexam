package token

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/redis.v5"
)

const TOKENKEY = "CLOUD:TOKEN:"

type AuthCodePayload struct {
	UID   string `json:"UID"`
	AppID string `json:"AppID"`
	Token string `json:"Token"`
}

//Response 客户端应答
type AuthResponse struct {
	//Code 错误码
	Code    int      `json:"Code"`
	Message string   `json:"Message"`
	Data    AuthResp `json:"Data"`
}
type AuthResp struct {
	AuthCode string
	Result   bool
}

//uid 用户ID
//accesstoken 酷派token
//appid 应用ID
// accounturl 账号acountapi服务获取AuthCode的url
//tm redis保存时间单位小时
func GetCloudToken(uid, accesstoken, appid, accounturl string, tm int, cache *redis.Client) (string, string, error) {

	var (
		authCode AuthCodePayload
		authcode string
		token    string
	)
	authCode.UID = uid
	authCode.AppID = appid
	authCode.Token = accesstoken
	data, err := json.Marshal(authCode)
	signMap := make(map[string]interface{})
	timestamp := time.Now().Unix()
	timestamps := strconv.FormatInt(timestamp, 10)
	signMap["Timestamp"] = timestamps
	signMap["Body"] = string(data)
	beego.Debug(signMap)
	sign, _ := GenSign(signMap, "429429dbb39c48c299686318bcf5e8e9")
	var header = make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("AppID", authCode.AppID)
	header.Set("Sign", sign)
	header.Set("Timestamp", timestamps)

	//	userurl := conf.Conf.Userurl + "/V1/Sync/AuthCodeForWeb"
	body, err := HttpPost(accounturl, string(data), header)
	if err != nil {
		beego.Error(err)
		return authcode, token, err

	}
	beego.Debug(string(body))
	var respbody AuthResponse

	err = json.Unmarshal(body, &respbody)
	if err != nil {
		// handle error
		beego.Error(err)
		return authcode, token, err
	}
	if respbody.Code == 0 {
		authcode = respbody.Data.AuthCode
		if err, token = CreateToken(uid, tm, cache); err != nil {
			return authcode, token, err
		}
	}
	return authcode, token, nil
}

func CreateToken(uid string, tm int, cache *redis.Client) (error, string) {
	h := md5.New()
	timestamp := time.Now().Unix()
	code := fmt.Sprintf("%v-%v", uid, timestamp)
	h.Write([]byte(code))
	token := hex.EncodeToString(h.Sum(nil))
	key := TOKENKEY + uid
	keepTime := time.Duration(tm) * time.Hour
	//存入redis
	err := cache.Set(key, token, keepTime).Err()
	beego.Debug(cache.Get(key).Result())
	//判断token是否失效
	if err != nil {
		return err, ""
	}
	return nil, token
}

func CheckToken(uid, token string, cache *redis.Client) bool {
	key := TOKENKEY + uid
	if user, err := cache.Get(key).Result(); err == nil {
		if user == token {
			return true
		}
	}
	return false
}

func GenSign(srcMap map[string]interface{}, secret string) (string, error) {
	//TODO GET Secret from redis
	srcMap["AppKey"] = secret
	keys := make([]string, len(srcMap))
	i := 0
	for k, _ := range srcMap {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	signStr := ""
	var value interface{}
	var v string
	for _, k := range keys {
		if k == "Sign" {
			continue
		}
		value = srcMap[k]
		switch value := value.(type) {

		case bool:
			v = strconv.FormatBool(value)
		case int:
			v = strconv.Itoa(value)
		case uint:
			v = strconv.FormatUint(uint64(value), 10)
		case int8:
			v = strconv.Itoa(int(value))
		case int16:
			v = strconv.Itoa(int(value))
		case int32:
			v = strconv.Itoa(int(value))
		case int64:
			v = strconv.FormatInt(value, 10)
		case float32:
			v = strconv.FormatFloat(float64(value), 'f', -1, 64)
		case float64:
			v = strconv.FormatFloat(value, 'f', -1, 64)
		case string:
			v = value
		default:
			continue

		}
		signStr = signStr + k + v
	}
	fmt.Println("Sign str is :", signStr)
	t := sha1.New()
	io.WriteString(t, signStr)
	out := fmt.Sprintf("%x", t.Sum(nil))
	return out, nil
}

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
