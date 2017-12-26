package services

import (
	"encoding/json"
	"notepad-api/conf"
	"notepad-api/model"
	"notepad-api/utils"

	"github.com/astaxie/beego"
)

//GetShareNote get one shared note's content
func GetShareNote(Ob map[string]interface{}) *utils.Response {
	var (
		share    *model.Share
		path     string
		answer   string
		content  string
		err      error
		attlist  []interface{}
		voteresp *utils.VoteResp
	)
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["Answer"].(string); ok {
		answer = Ob["Answer"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	share = model.NewShare(path)
	shareNote, err := share.Get()
	if shareNote == nil || err != nil {
		return utils.SHARE_ID_NOT_EXIST
	}
	if len(answer) == 0 {
		answer = "test"
	}
	if shareNote.Type == 1 {
		if shareNote.Answer != answer && shareNote.Password != answer {
			return utils.SHARE_ANSWER_NOT_TRUE
		}
	}

	noteID := shareNote.NoteID
	if content, err = utils.Cache.Get(noteID).Result(); err == nil && len(content) != 0 {

	} else {
		beego.Error(err)
		note := model.NewNote(noteID)
		notebook, err1 := note.GetWithPath()
		if notebook == nil || err1 != nil {
			return utils.NOTE_NOT_EXIST
		}
		content = notebook.Content
	}
	key := utils.SHAREKEY + noteID
	utils.ShareCache(key, content, 3600)
	if votes, err := utils.Cache.Get(path).Result(); err == nil && len(votes) != 0 {
		json.Unmarshal([]byte(votes), &voteresp)
	} else {
		voteresp = GetVote(noteID)
	}
	//	utils.Cache.Put(path, voteresp, 3600)
	filed := []string{"OpenNum"}
	OpenNum := shareNote.OpenNum + 1
	shareNote.OpenNum = OpenNum
	if err := shareNote.Update(filed); err != nil {
		beego.Error(err)
	}
	attach := model.NewAttachmentList(noteID)
	attachList, err := attach.GetList()
	if err != nil {
		beego.Error(err)
	}
	cloudService := GetCloudBoxService()
	for i := 0; i < len(attachList); i++ {
		var attachment utils.Attachment
		if attachList[i].Status == model.STATUS_DELETED {
			continue
		}
		attachment.Path = attachList[i].OK
		attachment.Url = cloudService.GetUrl(attachList[i].OK)
		attachment.Uuid = attachList[i].ID
		attachment.Type = attachList[i].Type
		attachment.ObType = attachList[i].OT
		attlist = append(attlist, attachment)
	}
	t := struct {
		Content        string
		Vote           interface{}
		AttachmentList interface{}
	}{content, voteresp, attlist}

	return utils.NewResponse(0, "", t)
}

//AddShare get share note url
func AddShare(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		share    *model.Share
		data     = make(map[string]interface{}, 0)
		note     *model.Note
		ID, path string
		Type     uint8
		file     []string
		isExist  bool
		//		Usn       int
		issuccess bool
	)
	share = model.NewShare("")
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
		note = model.NewNote(path)
		note.UID = userId
		notebook, err := note.Get()
		if notebook == nil || err != nil {
			return utils.NOTE_NOT_EXIST
		}
		share.NoteID = path
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["Type"].(float64); ok {
		Type = uint8(Ob["Type"].(float64))
		if Type != 1 && Type != 8 {
			if Type != 2 && Type != 4 {
				return utils.SHARE_TYPE_NOT_EXIST
			}
			share.Type = 1
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["Data"].(map[string]interface{}); ok {
		data = Ob["Data"].(map[string]interface{})
	}
	if sharenote, err := share.GetByNoteIDAndType(); sharenote != nil && err == nil {
		share.ShareID = sharenote.ShareID
		file = append(file, "CreateTime")
		share.CreateTime = utils.NowSecond()
		isExist = true

	} else {
		for i := 0; i < 10; i++ {
			ID = utils.CreateRandId(10)
			share = model.NewShare(ID)
			if !share.IsExist() {
				issuccess = true
				break
			}
		}
		if !issuccess {
			return utils.SHARE_WRONG
		}
		share.NoteID = path
		if Type != 1 && Type != 8 {
			share.Type = 1
		}
	}
	if Type == model.QUESTION_SHARE {
		if _, ok := data["Question"].(string); ok {
			share.Question = data["Question"].(string)
			if len([]rune(share.Question)) == 0 || len([]rune(share.Question)) > 30 {
				return utils.SHARE_QUESTION_LIMIT
			}
		} else {
			return utils.PARAMSFORMATERROR
		}
		if _, ok := data["Answer"].(string); ok {
			share.Answer = data["Answer"].(string)
			if len([]rune(share.Answer)) == 0 || len([]rune(share.Answer)) > 20 {
				return utils.SHARE_QUESTION_LIMIT
			}
		} else {
			return utils.PARAMSFORMATERROR
		}
	} else if Type == model.PASSWORD_SHARE {
		if _, ok := data["Password"].(string); ok {
			share.Password = data["Password"].(string)
			if len(share.Password) == 0 {
				return utils.SHARE_PASSWORD_NOT_NULL
			}
			if len([]rune(share.Password)) < 4 && len([]rune(share.Password)) > 12 {
				return utils.SHARE_PASSWORD_LIMIT
			}
		} else {
			return utils.PARAMSFORMATERROR
		}
	}
	if Type == model.VOTE_SHARE {
		data["Path"] = path
		t := VoteAdd(data, userId)
		beego.Debug(t)
		if t.Code != 0 && t.Code != 1023 {
			return t
		}
	}
	if isExist {
		file = append(file, "Password")
		file = append(file, "Answer")
		file = append(file, "Question")
		file = append(file, "Type")
		if err := share.Update(file); err != nil {
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}
		if Type == model.PASSWORD_SHARE {
			note.SetPasswordShare(model.STATUS_DELETED)
		}
		//		Usn = IncrUsn(userId)
	} else {
		if err := share.Add(); err != nil {
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}
		if Type == model.PASSWORD_SHARE {
			note.SetPasswordShare(model.STATUS_DELETED)
		}
		//		Usn = IncrUsn(userId)

	}
	//	AddNoteUsn(note.Path, userId, Usn)
	t := struct {
		Url string
	}{conf.Conf.ShareUrl + share.ShareID}
	return utils.NewResponse(0, "", t)
}

//GetShareMessage get all share message
func GetShareMessage(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		share *model.Share
		path  string
	)
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	note := model.NewNote(path)
	note.UID = userId
	notebook, err1 := note.Get()
	if notebook == nil || err1 != nil {
		return utils.NOTE_NOT_EXIST
	}
	share = model.NewShare("")
	share.NoteID = path
	sharenote, err := share.GetByNoteID()
	if sharenote == nil || err != nil {
		return utils.NOTE_NOT_SHARE
	}
	t := struct {
		Question string
		Answer   string
		Password string
	}{sharenote.Question, sharenote.Answer, sharenote.Password}
	return utils.NewResponse(0, "", t)
}

//GetQuestion get one note share question
func GetQuestion(Ob map[string]interface{}) *utils.Response {
	var share *model.Share
	var path string
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	share = model.NewShare(path)
	shareNote, err := share.Get()
	if shareNote == nil || err != nil {
		return utils.SHARE_ID_NOT_EXIST
	}
	noteID := shareNote.NoteID

	note := model.NewNote(noteID)
	notebook, err1 := note.GetWithPath()
	if notebook == nil || err1 != nil {
		return utils.NOTE_NOT_EXIST
	}
	t := struct {
		Question string
		Title    string
		Type     uint8
	}{shareNote.Question, notebook.Title, shareNote.Type}
	return utils.NewResponse(0, "", t)
}

//GetShareUrl get share note url
func GetShareUrl(Ob map[string]interface{}, userId string) *utils.Response {
	var share *model.Share
	var noteID string
	if _, ok := Ob["Path"].(string); ok {
		noteID = Ob["Path"].(string)
		note := model.NewNote(noteID)
		note.UID = userId
		if !note.IsExist() {
			return utils.NOTE_NOT_EXIST
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	share = model.NewShare("")
	share.NoteID = noteID
	if sharenote, err := share.GetByNoteID(); sharenote == nil || err != nil {
		return utils.NOTE_NOT_SHARE
	}
	shareUrl := beego.AppConfig.String("sharUrl")
	t := struct {
		Url string
	}{shareUrl + share.ShareID}
	return utils.NewResponse(0, "", t)
}

//CancleShare cancle one note share
func CancleShare(Ob map[string]interface{}, userId string) *utils.Response {
	var share *model.Share
	var noteID string
	if _, ok := Ob["Path"].(string); ok {
		noteID = Ob["Path"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	share = model.NewShare("")
	share.NoteID = noteID
	sharenote, err := share.GetByNoteID()
	if sharenote == nil || err != nil {
		return utils.NOTE_NOT_SHARE
	}
	err = sharenote.CancleShare()
	if err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	note := model.NewNote(sharenote.NoteID)
	note.UID = userId
	err = note.SetShare(model.NOTE_NOT_SHARE)
	if err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		IsSuccess bool
	}{true}
	return utils.NewResponse(0, "", t)
}
