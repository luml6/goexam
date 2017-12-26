package services

import (
	"notepad-api/model"
	"notepad-api/utils"
	"strconv"
	"strings"
	//	"strings"
	"time"

	"github.com/astaxie/beego"
)

//Create create new note
func Create(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		CategoryPath, CategoryName string
		ret                        *utils.Response
		note                       *model.Note
		filed                      []string
		Usn                        int
		attachlist                 []interface{}
		path                       string
	)
	path = utils.CreateId()
	note = model.NewNote(path)
	if _, ok := Ob["CategoryPath"].(string); ok {
		CategoryPath = Ob["CategoryPath"].(string)
		if CategoryPath != "" {
			cate := model.NewCategory(userId, CategoryPath)
			if !cate.IsExist() {
				return utils.CATEGORY_NOT_EXIST
			}
			category, _ := cate.Get()
			CategoryName = category.Name
			filed = append(filed, "CategoryPath")
			filed = append(filed, "CategoryName")
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["Title"].(string); ok {
		note.Title = Ob["Title"].(string)
		note.CategoryPath = CategoryPath
		note.CategoryName = CategoryName
		filed = append(filed, "Title")
		if len([]rune(note.Title)) > 80 {
			return utils.NOTE_TITLE_LIMIT
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["Content"].(string); ok {
		note.Content = Ob["Content"].(string)
		filed = append(filed, "Content")
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["CreateTime"]; ok {
		note.CreateTime = utils.GetTime(Ob["CreateTime"])
	} else {
		note.CreateTime = utils.NowSecond()
	}
	//	note.Source = 2
	filed = append(filed, "CreateTime")
	filed = append(filed, "UpdateTime")
	note.UpdateTime = note.CreateTime
	if _, ok := Ob["Author"].(string); ok {
		note.Author = Ob["Author"].(string)
		filed = append(filed, "Author")
	}
	if _, ok := Ob["Summary"].(string); ok {
		note.Summary = Ob["Summary"].(string)
		filed = append(filed, "Summary")
	}
	if _, ok := Ob["CheckSum"].(string); ok {
		note.CheckSum = Ob["CheckSum"].(string)
		filed = append(filed, "CheckSum")
	}
	if _, ok := Ob["IsPinned"].(float64); ok {
		if IsPinned, ok := Ob["IsPinned"].(float64); ok {
			note.IsPinned = uint8(IsPinned)
		} else {
			note.IsPinned = model.NOTE_NOT_PINNED
		}
		filed = append(filed, "IsPinned")
	}
	if _, ok := Ob["Attachlist"].([]interface{}); ok {
		attachlist = Ob["Attachlist"].([]interface{})
	}
	note.UID = userId
	Usn = IncrUsn(userId)
	filed = append(filed, "UID")
	filed = append(filed, "Status")
	filed = append(filed, "Usn")
	note.Usn = Usn
	if note.IsExist() {
		return utils.NOTE_ALREADY_EXIST
	}
	if err := note.Add(); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	var fileds = []string{"NotePath", "Status"}
	for i := 0; i < len(attachlist); i++ {
		attchID, _ := attachlist[i].(string)
		attch := model.NewAttachment(attchID)
		attch.Status = model.STATUS_PRECREATE
		attach, err := attch.Get()
		if err != nil || attach == nil {
			beego.Error(err)
			continue
		}
		attach.NotePath = path
		attach.Status = model.STATUS_NORMAL
		if err := attach.Update(fileds); err != nil {
			beego.Error(err)
		}
	}
	tm := time.Now().Unix()
	t := struct {
		Path string
		Usn  int
		Time int64
	}{note.Path, Usn, tm}
	ret = utils.NewResponse(0, "", t)
	//	}
	return ret
}

//DeleteByV2 delete one note
func DeleteByV2(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		note *model.Note
		path string
		Usn  int
	)
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
		note = model.NewNote(path)
		note.UID = userId
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["UpdateTime"]; ok {
		note.UpdateTime = utils.GetTime(Ob["UpdateTime"])
	} else {
		note.UpdateTime = utils.NowSecond()
	}
	key := utils.NOTELOCKKEY + path
	if utils.CheckLock(key) {
		return utils.NOTE_IS_LOCK
	}
	if !note.IsExist() {
		Usn = IncrUsn(userId)
		AddNoteUsn(note.Path, userId, Usn)
		tm := time.Now().Unix()
		t := struct {
			Path string
			Usn  int
			Time int64
		}{note.Path, Usn, tm}
		return utils.NewResponse(0, "", t)
	}
	if err := note.Delete(); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	sharekey := utils.SHAREKEY + note.Path
	//添加缓存
	if content, err := utils.Cache.Get(sharekey).Result(); err == nil && len(content) != 0 {

		if err := utils.Cache.Del(sharekey).Err(); err != nil {
			beego.Error(err)
		}
	}
	Usn = IncrUsn(userId)
	AddNoteUsn(note.Path, userId, Usn)

	tm := time.Now().Unix()
	t := struct {
		Path string
		Usn  int
		Time int64
	}{note.Path, Usn, tm}
	return utils.NewResponse(0, "", t)
}

//Delete delete one note
func Delete(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		note *model.Note
		path string
		Usn  int
	)
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
		note = model.NewNote(path)
		note.UID = userId
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["UpdateTime"]; ok {
		note.UpdateTime = utils.GetTime(Ob["UpdateTime"])
	} else {
		note.UpdateTime = utils.NowSecond()
	}
	if !note.IsExist() {
		Usn = IncrUsn(userId)
		AddNoteUsn(note.Path, userId, Usn)
		tm := time.Now().Unix()
		t := struct {
			Path string
			Usn  int
			Time int64
		}{note.Path, Usn, tm}
		return utils.NewResponse(0, "", t)
	}
	if err := note.Delete(); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	sharekey := utils.SHAREKEY + note.Path
	//添加缓存
	if content, err := utils.Cache.Get(sharekey).Result(); err == nil && len(content) != 0 {
		if err := utils.Cache.Del(sharekey).Err(); err != nil {
			beego.Error(err)
		}
	}
	Usn = IncrUsn(userId)
	AddNoteUsn(note.Path, userId, Usn)

	tm := time.Now().Unix()
	t := struct {
		Path string
		Usn  int
		Time int64
	}{note.Path, Usn, tm}
	return utils.NewResponse(0, "", t)
}

//DeleteAttach delete one note's attach
func DeleteAttach(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		attach   *model.Attachment
		path     []interface{}
		notePath string
		attlist  []string
	)
	if _, ok := Ob["Uuid"].([]interface{}); ok {
		path = Ob["Uuid"].([]interface{})
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["NotePath"].(string); ok {
		notePath = Ob["NotePath"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	attachone := model.NewAttachmentList(notePath)
	attachList, err := attachone.GetList()
	if err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	for i := 0; i < len(attachList); i++ {
		attlist = append(attlist, attachList[i].ID)
	}
	if len(path) == 0 {
		if err := DeleteAttachment(notePath); err != nil {
			beego.Error(err)
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}
	}
	for i := 0; i < len(path); i++ {
		id := path[i].(string)
		attach = model.NewAttachment(id)
		attach.NotePath = notePath
		if !attach.IsExist() {
			return utils.ATTACH_NOT_EXIST
		}
		if utils.CheckString(attlist, id) {
			continue
		}
		if err := attach.Remove(); err != nil {
			beego.Error(err)
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}
		cloudBoxService := GetCloudBoxService()
		cloudBoxService.DeleteQiuniu(id)

	}

	t := struct {
		Result string
	}{"success"}
	return utils.NewResponse(0, "", t)
}

//GetNote get one note message
func GetNote(Ob map[string]interface{}, userId string) *utils.Response {
	var path string
	var reNote utils.Note
	var status uint8 = 0
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["Status"].(float64); ok {
		status = uint8(Ob["Status"].(float64))
	}
	note := model.NewNote(path)
	note.UID = userId
	beego.Debug(note)
	notebook, err := note.GetByPathAndStatus(status)
	if notebook == nil || err != nil {
		return utils.NOTE_NOT_EXIST
	}
	reNote = NoteResp(*notebook)
	return utils.NewResponse(0, "", reNote)
}

//GetByUsn get one note message
//func GetByUsn(Ob map[string]interface{}, userId string) *utils.Response {
//	var path string
//	var reNote utils.Note
//	//	var status uint8 = 0
//	var usn int
//	if _, ok := Ob["Path"].(string); ok {
//		path = Ob["Path"].(string)
//	} else {
//		return utils.WITHOUT_PARAMETERS
//	}
//	if _, ok := Ob["Usn"].(float64); ok {
//		usn = int(Ob["Usn"].(float64))
//	}
//	note := model.NewNote(path)
//	note.UID = userId
//	//	note.Usn = usn
//	beego.Debug(note)
//	notebook, err := note.GetAllStatus()
//	if notebook == nil || err != nil {
//		return utils.NOTE_NOT_EXIST
//	}
//	if notebook.Usn == usn {
//		return utils.NOTE_NOT_CHANGE
//	}
//	reNote = NoteResp(*notebook)
//	return utils.NewResponse(0, "", reNote)
//}
func NoteResp(notebook model.Note) utils.Note {
	var reNote utils.Note
	reNote.Author = notebook.Author
	reNote.CategoryName = notebook.CategoryName
	reNote.CategoryPath = notebook.CategoryPath
	reNote.CheckSum = notebook.CheckSum
	reNote.Content = notebook.Content
	reNote.CreateTime = notebook.CreateTime.Unix()
	reNote.Summary = notebook.Summary
	reNote.Title = notebook.Title
	reNote.Usn = notebook.Usn
	reNote.IsPassworded = notebook.IsPassword
	reNote.UpdateTime = notebook.UpdateTime.Unix()
	reNote.IsPinned = notebook.IsPinned
	reNote.IsVote = notebook.IsVote
	reNote.Vote = GetVote(notebook.Path)
	//	reNote.IsPassworded = notebook.IsPassword
	attach := model.NewAttachmentList(notebook.Path)
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
		reNote.AttachmentList = append(reNote.AttachmentList, attachment)
	}
	return reNote
}

//Move move one note to other category
func Move(Ob map[string]interface{}, userId string) *utils.Response {
	var noteId, categoryId string
	filed := []string{"CategoryPath", "CategoryName"}
	var note *model.Note
	if _, ok := Ob["Path"].(string); ok {
		noteId = Ob["Path"].(string)
		note = model.NewNote(noteId)
		note.UID = userId
		if !note.IsExist() {
			return utils.NOTE_NOT_EXIST
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["CategoryPath"].(string); ok {
		categoryId = Ob["CategoryPath"].(string)
		if categoryId != "" {
			cate := model.NewCategory(userId, categoryId)
			category, err := cate.Get()
			if category != nil && err == nil {
				note.CategoryName = category.Name
				note.CategoryPath = category.Path
			} else {
				beego.Error(err)
				return utils.CATEGORY_NOT_EXIST
			}
		}

	} else {
		return utils.WITHOUT_PARAMETERS
	}
	//	note.UpdateTime = utils.NowSecond()
	if err := note.Update(filed); err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	Usn := IncrUsn(userId)
	AddNoteUsn(note.Path, userId, Usn)
	tm := time.Now().Unix()
	t := struct {
		Path string
		Usn  int
		Time int64
	}{note.Path, Usn, tm}
	return utils.NewResponse(0, "", t)
}

func UpdateByV2(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		note       *model.Note
		noteId     string
		filed      []string
		attachlist []interface{}
		Usn        int
		err        error
		attchtype  utils.AttchType
		islock     bool
		istrue     bool
		createID   string
		checkSum   string
	)
	tm := time.Now().Unix()
	if _, ok := Ob["Path"].(string); ok {
		noteId = Ob["Path"].(string)
		note = model.NewNote(noteId)
		note.UID = userId
		note, _ = note.GetAllStatus()
		if note == nil {
			return utils.NOTE_NOT_EXIST
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	key := utils.NOTELOCKKEY + noteId
	if utils.CheckLock(key) {
		return utils.NOTE_IS_LOCK
	}
	if _, ok := Ob["CheckSum"].(string); ok {
		checkSum = Ob["CheckSum"].(string)
		filed = append(filed, "CheckSum")
	} else {
		istrue = true
	}
	if _, ok := Ob["Usn"].(float64); ok {
		Usn = int(Ob["Usn"].(float64))
		if note.Usn != Usn && note.Status == model.STATUS_NORMAL && Usn != 0 && !istrue {
			utils.LockCache(key, 1, 30)
			islock = true
			filed = append(filed, "ConflictCount")
			conflict := model.NewConflict(utils.CreateId())
			conflict.FatherPath = noteId
			conflict.CheckSum = checkSum
			if conflict, err := conflict.GetBySum(); err == nil && conflict != nil {
				t := struct {
					Path       string
					Usn        int
					Time       int64
					AttachType interface{}
				}{note.Path, note.Usn, tm, nil}
				return utils.NewResponse(0, "", t)
			}
			if createID, err = Merge(*note, noteId, checkSum); err != nil {
				beego.Error(err)
				return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
			}
			note.ConflictCount += 1
		}
	}
	if !istrue {
		note.CheckSum = checkSum
	}
	if _, ok := Ob["CategoryPath"].(string); ok {
		categoryId := Ob["CategoryPath"].(string)
		if categoryId != "" {
			filed = append(filed, "CategoryPath", "CategoryName")
			cate := model.NewCategory(userId, categoryId)
			category, err := cate.Get()
			if category != nil && err == nil {
				note.CategoryName = category.Name
				note.CategoryPath = category.Path
			} else {
				beego.Error(err)
				//				return utils.CATEGORY_NOT_EXIST
			}
		}

	}
	if _, ok := Ob["Author"].(string); ok {
		note.Author = Ob["Author"].(string)
		filed = append(filed, "Author")
	}

	if _, ok := Ob["Title"].(string); ok {
		note.Title = Ob["Title"].(string)
		if len([]rune(note.Title)) > 80 {
			return utils.NOTE_TITLE_LIMIT
		}
		filed = append(filed, "Title")
	}
	if _, ok := Ob["Summary"].(string); ok {
		note.Summary = Ob["Summary"].(string)
		filed = append(filed, "Summary")
	}

	if _, ok := Ob["Content"].(string); ok {
		note.Content = Ob["Content"].(string)
		filed = append(filed, "Content")
	}
	if _, ok := Ob["UpdateTime"]; ok {
		note.UpdateTime = utils.GetTime(Ob["UpdateTime"])
	} else {
		note.UpdateTime = utils.NowSecond()
	}
	filed = append(filed, "UpdateTime")
	if _, ok := Ob["IsPinned"].(float64); ok {
		if IsPinned, ok := Ob["IsPinned"].(float64); ok {
			if uint8(IsPinned) != 255 {
				note.IsPinned = uint8(IsPinned)
				filed = append(filed, "IsPinned")
			}
		}
	}
	note.Status = model.STATUS_NORMAL
	filed = append(filed, "Status")
	if err := note.Update(filed); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	//添加缓存
	sharekey := utils.SHAREKEY + note.Path
	if content, err := utils.Cache.Get(sharekey).Result(); err == nil && len(content) != 0 {
		utils.ShareCache(sharekey, note.Content, 3600)
	}
	Usn = IncrUsn(userId)
	AddNoteUsn(note.Path, userId, Usn)

	if _, ok := Ob["Attachlist"].([]interface{}); ok {
		attachlist = Ob["Attachlist"].([]interface{})
		if attchtype, err = Attachlist(attachlist, note.Path); err != nil {
			beego.Error(err)
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}
	}
	if islock {
		utils.Cache.Del(key)
	}

	t := struct {
		Path       string
		Usn        int
		Time       int64
		NewPath    string `json:",omitempty"`
		AttachType interface{}
	}{note.Path, Usn, tm, createID, attchtype}
	return utils.NewResponse(0, "", t)
}

func Attachlist(attachlist []interface{}, notePath string) (utils.AttchType, error) {
	var existarray []string
	var attchtype utils.AttchType
	attachone := model.NewAttachmentList(notePath)
	attachList, err := attachone.GetList()
	if err != nil {
		beego.Error(err)
		return attchtype, err
	}
	for i := 0; i < len(attachlist); i++ {
		attachID := attachlist[i].(string)
		attach := model.NewAttachment(attachID)
		if att, err := attach.GetAllStatus(); err != nil || att == nil {
			beego.Error(err)
			return attchtype, err
		} else {
			//			if att.Status != model.STATUS_DELETED {
			existarray = append(existarray, attachID)
			if att.Status == model.STATUS_PRECREATE || att.Status == model.STATUS_DELETED {
				att.NotePath = notePath
				att.Status = model.STATUS_NORMAL
				fileds := []string{"NotePath", "Status"}
				beego.Debug(att)
				if err := att.Update(fileds); err != nil {
					beego.Error(err)
				}
			}
			//			}
		}
	}
	for i := 0; i < len(attachList); i++ {
		if !utils.CheckString(existarray, attachList[i].ID) {
			if err := attachList[i].SetStatus(model.STATUS_DELETED); err != nil {
				beego.Error(err)
				continue
			}
		} else {
			switch attachList[i].Type {
			case "1":
				attchtype.HasImg = true
			case "2":
				attchtype.HasMusic = true
			case "3":
				attchtype.HasVideo = true
			default:
			}

		}

	}
	return attchtype, nil
}

//处理冲突
func Merge(note model.Note, noteId, checkSum string) (string, error) {
	path := utils.CreateId()
	conflict := model.NewConflict(utils.CreateId())
	conflict.FatherPath = noteId
	conflict.SonPath = path
	conflict.CheckSum = checkSum
	note.Path = path
	note.Title = note.Title + "(" + strconv.Itoa(note.ConflictCount+1) + ")" + "(" + time.Now().Format("15:04:05") + ")"
	note.IsPassword = 0
	note.IsShare = 0
	note.IsVote = 0
	attachList := model.NewAttachmentList(noteId)
	attachs, _ := attachList.GetList()
	cloudService := GetCloudBoxService()
	for i := 0; i < len(attachs); i++ {
		attachs[i].NotePath = note.Path
		id := utils.CreateId()
		if err := cloudService.Upload(attachs[i].ID, id); err != nil {
			beego.Error(err)
			continue
		}
		note.Content = strings.Replace(note.Content, attachs[i].ID, id, -1)
		attachs[i].ID = id
		attachs[i].OK = id
		if err := attachs[i].Add(); err != nil {
			beego.Error(err)
			return path, err
		}
	}
	if err := note.Add(); err != nil {
		beego.Error(err)
		return path, err
	}
	//	conflict.CreateTime = time.Now().Format(model.FOTMAT_TIME_STRING)
	if err := conflict.Add(); err != nil {
		beego.Error(err)
	}
	return path, nil
}

//Update update one note
func Update(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		note       *model.Note
		noteId     string
		filed      []string
		attachlist []interface{}
		Usn        int
		attachtype utils.AttchType
		err        error
	)
	if _, ok := Ob["Path"].(string); ok {
		noteId = Ob["Path"].(string)
		note = model.NewNote(noteId)
		note.UID = userId
		note, _ = note.GetAllStatus()
		if note == nil {
			return utils.NOTE_NOT_EXIST
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["Author"].(string); ok {
		note.Author = Ob["Author"].(string)
		filed = append(filed, "Author")
	}

	if _, ok := Ob["Title"].(string); ok {
		note.Title = Ob["Title"].(string)
		if len([]rune(note.Title)) > 80 {
			return utils.NOTE_TITLE_LIMIT
		}
		filed = append(filed, "Title")
	}
	if _, ok := Ob["Summary"].(string); ok {
		note.Summary = Ob["Summary"].(string)
		filed = append(filed, "Summary")
	}
	if _, ok := Ob["CheckSum"].(string); ok {
		note.CheckSum = Ob["CheckSum"].(string)
		filed = append(filed, "CheckSum")
		filed = append(filed, "UpdateTime")
	}

	if _, ok := Ob["Content"].(string); ok {
		note.Content = Ob["Content"].(string)
		filed = append(filed, "Content")
	}
	if _, ok := Ob["UpdateTime"]; ok {
		note.UpdateTime = utils.GetTime(Ob["UpdateTime"])
	} else {
		note.UpdateTime = utils.NowSecond()
	}

	if _, ok := Ob["IsPinned"].(float64); ok {
		if IsPinned, ok := Ob["IsPinned"].(float64); ok {
			if uint8(IsPinned) != 255 {
				note.IsPinned = uint8(IsPinned)
				filed = append(filed, "IsPinned")
			}
		}
	}
	note.Status = model.STATUS_NORMAL
	filed = append(filed, "Status")
	if err := note.Update(filed); err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	key := utils.SHAREKEY + note.Path
	//添加缓存
	if content, err := utils.Cache.Get(key).Result(); err == nil && len(content) != 0 {
		utils.ShareCache(key, note.Content, 3600)
	}
	Usn = IncrUsn(userId)
	AddNoteUsn(note.Path, userId, Usn)

	if _, ok := Ob["Attachlist"].([]interface{}); ok {
		attachlist = Ob["Attachlist"].([]interface{})
		if attachtype, err = Attachlist(attachlist, note.Path); err != nil {
			beego.Error(err)
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}
	}
	tm := time.Now().Unix()
	t := struct {
		Path       string
		Usn        int
		Time       int64
		AttachType interface{}
	}{note.Path, Usn, tm, attachtype}
	return utils.NewResponse(0, "", t)
}

//GetOtherNoteList get other category all note message
func GetOtherNoteList(userId string) *utils.Response {
	noteList := model.NewNoteList("", userId)
	notes, err := noteList.GetListNoCategory()
	if err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	return utils.NewResponse(0, "", convertNote(notes))

}

//GetNoteID create one null note
func GetNoteID(userId string) *utils.Response {
	var note model.Note
	note.UID = userId
	path := utils.CreateId()
	note.Path = path
	note.Status = model.STATUS_PRECREATE
	if err := note.Create(); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	//	IncrUsn(userId)
	t := struct {
		Path string
	}{note.Path}

	return utils.NewResponse(0, "", t)
}

//GetAllNoteList get other category all note message
func GetAllNoteList(ob map[string]interface{}, userId string) interface{} {
	var (
		path   string
		result interface{}
	)
	if _, ok := ob["Path"].(string); ok {
		path = ob["Path"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}

	switch path {
	case "all":
		result = GetAllnote(userId)
	case "recycle":
		result = GetRecycleBin(userId)
	default:
		result = Getnotes(path, userId, false)
	}
	return result

}

//AddNoteUsn add one note usn
func AddNoteUsn(path, userId string, Usn int) {
	note := model.NewNote(path)
	note.UID = userId
	filed := []string{"Usn"}
	note.Usn = Usn
	beego.Debug("inc Usn")
	if err := note.Update(filed); err != nil {
		beego.Error(err)
	}
	return
}

//RecoverNote recover one deleted note
func RecoverNote(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		note *model.Note
		path string
		Usn  int
	)
	filed := []string{"Status", "UpdateTime", "CategoryPath", "CategoryName", "IsPinned"}

	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
		note = model.NewNote(path)
		note.UID = userId
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	notebook, err := note.GetAllStatus()
	if err != nil || notebook == nil {
		return utils.NOTE_NOT_EXIST
	}
	categoryPath := notebook.CategoryPath
	cate := model.NewCategory(userId, categoryPath)
	if cate.IsExist() {
		if err = note.Recover(); err != nil {
			beego.Error(err)
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}
	} else {
		note.UpdateTime = utils.NowSecond()
		note.CategoryName = ""
		note.CategoryPath = ""
		note.IsPinned = 0
		note.Status = model.STATUS_NORMAL
		if err = note.Update(filed); err != nil {
			beego.Error(err)
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}
	}
	Usn = IncrUsn(userId)
	AddNoteUsn(note.Path, userId, Usn)
	tm := time.Now().Unix()
	t := struct {
		Path string
		Usn  int
		Time int64
	}{note.Path, Usn, tm}
	return utils.NewResponse(0, "", t)

}

//RemoveNote remove one note from recycle bin
func RemoveNote(Ob map[string]interface{}, userId string) *utils.Response {
	var note *model.Note
	var path string
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
		note = model.NewNote(path)
		note.UID = userId
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	notebook, err := note.GetAllStatus()
	if err != nil || notebook == nil {
		beego.Error(err)
		return utils.NOTE_NOT_EXIST
	}
	if err = note.TrashClear(); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	DeleteAttachment(note.Path)
	t := struct {
		Path string
	}{path}
	return utils.NewResponse(0, "", t)

}
