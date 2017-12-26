package services

//七牛云存储上传下载凭证管理

import (
	"errors"
	"fmt"
	//	"io/ioutil"
	"notepad-api/conf"
	"notepad-api/model"
	"notepad-api/utils"
	"strings"

	"github.com/astaxie/beego"
	qiniuConf "qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
)

type CloudBoxService struct {
	maxUploadExpire   int
	maxDownloadExpire int
	noteBucket        string
	publicBucket      string
	accessKey         string
	secretKey         string
	domain            string
	callbackurl       string
}

var cloudService *CloudBoxService

//初始化七牛云存储配置
func init_service() {
	if cloudService == nil {
		cloudService = new(CloudBoxService)
	}

	cloudService.accessKey = beego.AppConfig.String("qiniu::access_key")
	cloudService.secretKey = beego.AppConfig.String("qiniu::secret_key")
	cloudService.noteBucket = beego.AppConfig.String("qiniu::note_bucket")
	cloudService.publicBucket = beego.AppConfig.String("qiniu::public_bucket")
	cloudService.domain = beego.AppConfig.String("qiniu::domain")
	cloudService.callbackurl = beego.AppConfig.String("qiniu::callbackurl")

	if num, err := beego.AppConfig.Int("qiniu::max_upload_expire"); err != nil {
		cloudService.maxUploadExpire = 3600
	} else {
		cloudService.maxUploadExpire = num
	}

	if num, err := beego.AppConfig.Int("qiniu::max_download_expire"); err != nil {
		cloudService.maxDownloadExpire = 3600
	} else {
		cloudService.maxDownloadExpire = num
	}

	qiniuConf.ACCESS_KEY = cloudService.accessKey
	qiniuConf.SECRET_KEY = cloudService.secretKey
}

func init() {
	init_service()
}

//获取云存储服务
func GetCloudBoxService() *CloudBoxService {
	if cloudService == nil {
		init_service()
	}

	return cloudService
}

func paraGetUint32(field string, args map[string]interface{}) (uint32, error) {
	var ret uint32
	val, ok := args[field]
	if !ok {
		return ret, errors.New(field + " not found")
	}

	if tmp, ok := val.(float64); !ok {
		return ret, errors.New(field + " wrong type")
	} else {
		ret = uint32(tmp)
	}

	return ret, nil
}

func paraGetString(field string, args map[string]interface{}) (string, error) {
	var ret string
	val, ok := args[field]
	if !ok {
		return ret, errors.New(field + " not found")
	}

	if tmp, ok := val.(string); !ok {
		return ret, errors.New(field + " wrong type")
	} else {
		ret = tmp
	}

	return ret, nil
}

func isPublicKey(key string) bool {
	return strings.HasPrefix(key, utils.ATTACH_PUBLIC_PREFIX)
}

func (serv *CloudBoxService) Upload(fileKey, key string) error {

	// new一个Bucket对象
	c := kodo.New(0, nil)
	p := c.Bucket(serv.noteBucket)
	baseUrl := kodo.MakeBaseUrl(serv.domain, fileKey)
	policy := kodo.GetPolicy{}
	//调用MakePrivateUrl方法返回url
	fullURL := c.MakePrivateUrl(baseUrl, &policy)
	// 调用Fetch方法
	err := p.Fetch(nil, key, fullURL)
	//	if err != nil {
	//		fmt.Println("bucket.Fetch failed:", err)
	//	} else {
	//		fmt.Println("fetch success")
	//	}
	return err
}

func (serv *CloudBoxService) GetUrl(key string) string {
	c := kodo.New(0, nil)
	baseUrl := kodo.MakeBaseUrl(serv.domain, key)
	policy := kodo.GetPolicy{}
	//调用MakePrivateUrl方法返回url
	fullURL := c.MakePrivateUrl(baseUrl, &policy)
	return fullURL
}

//获取上传凭证
func (serv *CloudBoxService) GetUploadToken(args map[string]interface{}) *utils.Response {
	var (
		fileKey string
		expires uint32
		aid     string
		noteID  string
		err     error
	)

	if fileKey, err = paraGetString("FileKey", args); err != nil {
		return utils.WITHOUT_PARAMETERS
	}

	if aid, err = paraGetString("AccountID", args); err != nil {
		return utils.WITHOUT_PARAMETERS
	}

	if noteID, err = paraGetString("NoteID", args); err != nil {
		return utils.WITHOUT_PARAMETERS
	}

	if expires, err = paraGetUint32("Expires", args); err != nil {
		return utils.WITHOUT_PARAMETERS
	}

	note := model.NewNote(noteID)
	note.UID = aid
	if note, err = note.Get(); note == nil {
		return utils.NOTE_NOT_EXIST
	} else if err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}

	c := kodo.New(0, nil)

	//设置上传的策略
	bucket := serv.noteBucket

	policy := &kodo.PutPolicy{
		Scope: bucket + "/" + fileKey,
		//设置Token过期时间
		Expires: expires,
	}

	//生成一个上传token
	token := c.MakeUptoken(policy)

	ret := make(map[string]interface{}, 3)
	ret["UploadToken"] = token
	ret["Expires"] = expires
	ret["Domain"] = serv.domain

	return utils.NewResponse(utils.SUCCESS_CODE, "", ret)
}

//获取下载url
func (serv *CloudBoxService) GetDownloadURL(args map[string]interface{}) *utils.Response {
	var (
		fileKey string
		expires uint32
		aid     string
		err     error
		shareId string
		types   uint32
	)
	types, _ = paraGetUint32("Type", args)
	if fileKey, err = paraGetString("FileKey", args); err != nil {
		return utils.WITHOUT_PARAMETERS
	}
	if shareId, err = paraGetString("ShareID", args); err != nil {
		if aid, err = paraGetString("AccountID", args); err != nil {
			return utils.WITHOUT_PARAMETERS
		}
	}

	if expires, err = paraGetUint32("Expires", args); err != nil {
		return utils.WITHOUT_PARAMETERS
	}

	att := model.NewAttachment(fileKey)
	att.Status = model.STATUS_NORMAL
	attpub := model.NewPublicAttach(fileKey)
	attpub.Status = model.STATUS_NORMAL
	if att, err = att.Get(); att == nil {
		if attpub, err = attpub.Get(); attpub == nil {
			return utils.ATTACH_NOT_EXIST
		} else if err != nil {
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}

	} else if err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	if shareId != "" {
		resp := CheckNoteID(shareId, att.NotePath)
		if resp.Code != 0 {
			return resp
		}
	}
	beego.Debug(aid)

	baseUrl := kodo.MakeBaseUrl(serv.domain, fileKey)
	policy := kodo.GetPolicy{
		Expires: expires,
	}
	if types == 1 {
		baseUrl += conf.Conf.PhotoStyle
	}
	//生成一个client对象
	c := kodo.New(0, nil)
	//调用MakePrivateUrl方法返回url
	fullURL := c.MakePrivateUrl(baseUrl, &policy)

	if len(fullURL) == 0 {
		beego.Debug(baseUrl + ": get fullURL error")
		return utils.NewResponse(utils.SYSTEM_CODE, "", nil)
	}

	ret := make(map[string]interface{}, 4)
	ret["FileEncryption"] = att.FileEncryption
	ret["DownloadURL"] = fullURL
	ret["Expires"] = expires

	return utils.NewResponse(utils.SUCCESS_CODE, "", ret)
}

func (serv *CloudBoxService) GetUploadTokenWithCallback(args map[string]interface{}) *utils.Response {
	var (
		fileKey        string
		expires        uint32
		aid            string
		fileEncryption string
		fileType       string
		fileName       string
		err            error
	)

	if fileKey, err = paraGetString("FileKey", args); err != nil {
		return utils.WITHOUT_PARAMETERS
	}

	if aid, err = paraGetString("AccountID", args); err != nil {
		return utils.WITHOUT_PARAMETERS
	}

	if fileEncryption, err = paraGetString("FileEncryption", args); err != nil {
		return utils.WITHOUT_PARAMETERS
	}

	//	if noteID, err = paraGetString("NotePath", args); err != nil {

	//		return utils.WITHOUT_PARAMETERS
	//	}

	if expires, err = paraGetUint32("Expires", args); err != nil {
		return utils.WITHOUT_PARAMETERS
	}

	if fileType, err = paraGetString("FileType", args); err != nil {
		return utils.WITHOUT_PARAMETERS
	}

	if fileName, err = paraGetString("FileName", args); err != nil {
		return utils.WITHOUT_PARAMETERS
	}

	//	note := model.NewNote(noteID)
	//	note.UID = aid
	//	if note, err = note.Get(); note == nil {
	//		note = model.NewNote(noteID)
	//		note.UID = aid
	//		if note, err = note.GetByPathAndStatus(model.STATUS_PRECREATE); note == nil {
	//			return utils.NOTE_NOT_EXIST
	//		}
	//	} else if err != nil {
	//		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	//	}

	beego.Debug("start get uploadtoken")

	var callbackBodyFormat string = `{"Uuid":$(key),"Size": $(fsize),"Hash": $(etag), "FileType":"%s", "FileName":"%s","FileEncryption":"%s","Uid":"%s"}`
	callbackbody := fmt.Sprintf(callbackBodyFormat, fileType, fileName, fileEncryption, aid)
	c := kodo.New(0, nil)

	//设置上传的策略
	bucket := serv.noteBucket

	policy := &kodo.PutPolicy{
		Scope: bucket + "/" + fileKey,
		//设置Token过期时间
		Expires: expires,

		CallbackUrl:      serv.callbackurl,
		CallbackBody:     callbackbody,
		CallbackBodyType: "application/json",
	}

	//生成一个上传token
	token := c.MakeUptoken(policy)

	ret := make(map[string]interface{}, 3)
	ret["UploadToken"] = token
	ret["Expires"] = expires
	ret["Domain"] = serv.domain

	return utils.NewResponse(utils.SUCCESS_CODE, "", ret)
}

func CheckNoteID(shareID, NoteID string) *utils.Response {
	var err error
	share := model.NewShare(shareID)
	if share, err = share.Get(); err != nil || share == nil {
		return utils.SHARE_ID_NOT_EXIST
	}
	if share.NoteID != NoteID {
		return utils.PARAMSFORMATERROR
	}
	return utils.NewResponse(0, "", nil)
}

func (serv *CloudBoxService) DeleteQiuniu(key string) {
	c := kodo.New(0, nil)
	bucket := serv.noteBucket
	p := c.Bucket(bucket)

	//调用Delete方法删除文件
	res := p.Delete(nil, key)
	//打印返回值以及出错信息
	if res == nil {
		fmt.Println("Delete success")
	} else {
		fmt.Println(res)
	}
}
