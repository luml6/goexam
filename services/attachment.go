package services

import (
	"notepad-api/conf"
	"notepad-api/model"
	"notepad-api/utils"

	"github.com/astaxie/beego"
)

//Upload upload a new attachment
func Upload(noteId, uuid, fileType, fileName, UserId string) *utils.Response {
	attach := model.NewAttachment(uuid)
	note := model.NewNote(noteId)
	note.UID = UserId
	if !note.IsExist() {
		return utils.NOTE_NOT_EXIST
	}
	attach.NotePath = noteId
	attach.Name = fileName
	attach.Type = fileType
	url := conf.Conf.DownUrl + noteId
	attach.OK = url
	attach.OT = model.ATT_STORE_TYPE_ZEUSIS
	attach.UID = UserId
	if err := attach.Add(); err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}

	t := struct {
		Path string
	}{url}
	return utils.NewResponse(0, "", t)
}

//AddPublic add one new public  note attachment
func AddPublic(Ob map[string]interface{}, userId string) *utils.Response {
	var attach *model.PublicAttach
	if _, ok := Ob["Uuid"]; ok {
		path := Ob["Uuid"].(string)
		attach = model.NewPublicAttach(path)
		attach.Key = path
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["FileType"]; ok {
		attach.Type = Ob["FileType"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if attach.IsExist() {
		return utils.ATTACH_UUID_ALREADY_EXIST
	}
	if err := attach.Add(); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		Path string
	}{attach.ID}

	return utils.NewResponse(0, "", t)
}

//Add add one new note attachment
func Add(Ob map[string]interface{}, userId string) *utils.Response {
	var attach *model.Attachment
	if _, ok := Ob["Uuid"]; ok {
		path := Ob["Uuid"].(string)
		attach = model.NewAttachment(path)
		attach.OK = path
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["FileType"]; ok {
		attach.Type = Ob["FileType"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["FileName"]; ok {
		attach.Name = Ob["FileName"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["Size"]; ok {
		attach.Size = int(Ob["Size"].(float64))
	}
	//	attach.Source = 2
	serv := GetCloudBoxService()
	filekey := utils.GetQiniuhash(serv.noteBucket, attach.ID)
	if len(filekey) == 0 {
		return utils.ATTACH_NOT_EXIST
	}
	beego.Debug(filekey)
	attach.FileEncryption = filekey

	if _, ok := Ob["NotePath"]; ok {
		attach.NotePath = Ob["NotePath"].(string)
		note := model.NewNote(attach.NotePath)
		note.UID = userId
		if !note.IsExist() {
			return utils.NOTE_NOT_EXIST
		}
		attach.Status = model.STATUS_NORMAL
	} else {
		attach.Status = model.STATUS_PRECREATE
	}
	attach.UID = userId

	if attach.IsExist() {
		return utils.ATTACH_UUID_ALREADY_EXIST
	}
	attach.OT = model.ATT_STORE_TYPE_COULD
	if err := attach.Add(); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		Path string
	}{attach.ID}

	return utils.NewResponse(0, "", t)
}

//AttachCallback qiniu callback api
func AttachCallback(Ob map[string]interface{}) utils.CloudResp {

	var (
		attach *model.Attachment
		size   int
		hash   string
	)
	if _, ok := Ob["Uuid"].(string); ok {
		path := Ob["Uuid"].(string)
		attach = model.NewAttachment(path)
		attach.OK = path
		attach.Status = model.STATUS_PRECREATE
	} else {
		return setCallBack(false, "", "without paramters")
	}

	if _, ok := Ob["FileType"].(string); ok {
		attach.Type = Ob["FileType"].(string)
	} else {
		return setCallBack(false, attach.ID, "without paramters")
	}
	if _, ok := Ob["Uid"].(string); ok {
		attach.UID = Ob["Uid"].(string)
	} else {
		return setCallBack(false, attach.ID, "without paramters")
	}
	if _, ok := Ob["FileName"].(string); ok {
		attach.Name = Ob["FileName"].(string)
	} else {
		return setCallBack(false, attach.ID, "without paramters")
	}
	if _, ok := Ob["FileEncryption"]; ok {
		attach.FileEncryption = Ob["FileEncryption"].(string)
	} else {
		return setCallBack(false, attach.ID, "without paramters")
	}
	if _, ok := Ob["Size"]; ok {
		size = int(Ob["Size"].(float64))
	} else {
		return setCallBack(false, attach.ID, "without paramters")
	}
	if _, ok := Ob["Hash"]; ok {
		hash = Ob["Hash"].(string)
	} else {
		return setCallBack(false, attach.ID, "without paramters")
	}
	//	if _, ok := Ob["NoteId"].(string); ok {
	//		attach.NotePath = Ob["NoteId"].(string)
	//		note := model.NewNote(attach.NotePath)
	//		note.UID = attach.UID
	//		if !note.IsExist() {
	//			if note, _ = note.GetByPathAndStatus(model.STATUS_PRECREATE); note == nil {
	//				return setCallBack(false, attach.ID, "note not exist")
	//			}
	//		}
	//		note.Path = attach.NotePath
	//		note.UID = attach.UID
	//		notes, err := note.Get()
	//		if notes != nil && err == nil {
	//			attach.UID = notes.UID
	//		} else {
	//			beego.Debug(notes, err)
	//		}
	//	} else {
	//		return setCallBack(false, attach.ID, "without paramters")
	//	}

	if attach.IsExist() {
		return setCallBack(false, attach.ID, "attchment already exit")
	}
	if len(attach.FileEncryption) != 0 {
		if hash != attach.FileEncryption {
			beego.Debug("upload default")
			serv := GetCloudBoxService()
			serv.DeleteQiuniu(attach.ID)
			return setCallBack(false, attach.ID, "attchment upload default")
		}
	}
	attach.Size = size
	beego.Debug("end")
	attach.OT = model.ATT_STORE_TYPE_COULD
	if err := attach.Add(); err != nil {
		beego.Error(err)
		return setCallBack(false, attach.ID, err.Error())
	}
	return setCallBack(true, attach.ID, "")
}

func DeleteAttachment(noteId string) error {
	attachList := model.NewAttachmentList(noteId)
	attachs, _ := attachList.GetList()
	err := attachList.RemoveAll()
	if err == nil {
		if attachs != nil {
			for i := 0; i < len(attachs); i++ {
				//				if attachs[i].IsExist() {
				cloudBoxService := GetCloudBoxService()
				cloudBoxService.DeleteQiuniu(attachs[i].ID)
				//				}

			}
		}
	}
	return err
}
func setCallBack(success bool, key, msg string) utils.CloudResp {
	var t utils.CloudResp
	t.Key = key
	t.Message = msg
	t.Success = success
	return t
}

func Check(Ob map[string]interface{}) *utils.Response {
	var attach *model.Attachment
	if _, ok := Ob["Uuid"].(string); ok {
		path := Ob["Uuid"].(string)
		attach = model.NewAttachment(path)
	}
	if !attach.IsExist() {
		return utils.ATTACH_NOT_EXIST
	}
	t := struct {
		IsExist bool
	}{true}
	return utils.NewResponse(0, "", t)
}
