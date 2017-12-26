package services

import (
	"notepad-api/model"
	"notepad-api/utils"

	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

//Create create new category
func CreateCategory(Ob map[string]interface{}, userId string) *utils.Response {
	var cate *model.Category
	tm := time.Now().Unix()
	if _, ok := Ob["Name"].(string); ok {
		path := utils.CreateId()
		cate = model.NewCategory(userId, path)
		cate.Name = Ob["Name"].(string)
		if cate.Name == "" {
			return utils.CATEGORY_NAME_NOT_NULL
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["CreateTime"]; ok {
		cate.CreateTime = utils.GetTime(Ob["CreateTime"])
	} else {
		cate.CreateTime = utils.NowSecond()
	}
	//	cate.Source = 2
	cate.UpdateTime = cate.CreateTime
	category, err := cate.GetByNameAndUID(cate.Name, cate.UID)
	if category != nil || err != nil {
		beego.Error(utils.CATEGORY_NAME_EXIST)
		t := struct {
			Path string
			Usn  int
			Time int64
		}{category.Path, category.Usn, tm}
		return utils.NewResponse(0, "", t)

	}
	cates, _ := cate.GetRecycleByNameAndUID(cate.Name, cate.UID)
	if cates != nil {
		//		noteList := model.NewNoteList(cates.Path, userId)
		//		if err := noteList.DeleteCategory(); err != nil {
		//			beego.Error(err)
		//		}
		cates.Remove()
	}
	if cate.IsExist() {
		return utils.CATEGORY_NOT_EXIST
	}
	Usn := IncrUsn(userId)
	cate.Usn = Usn
	cate.UID = userId
	if err := cate.Add(); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}

	t := struct {
		Path string
		Usn  int
		Time int64
	}{cate.Path, Usn, tm}

	return utils.NewResponse(0, "", t)
}

//DeleteCategory delete one category
func DeleteCategory(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		cate *model.Category
		path string
		Usn  int
	)
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
		cate = model.NewCategory(userId, path)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["UpdateTime"]; ok {
		cate.UpdateTime = utils.GetTime(Ob["UpdateTime"])
	} else {
		cate.UpdateTime = utils.NowSecond()
	}
	if !cate.IsExist() {
		return utils.CATEGORY_NOT_EXIST
	}
	noteList := model.NewNoteList(path, userId)
	notelist, err := noteList.GetList()
	if err == nil {
		for i := 0; i < len(notelist); i++ {
			Usn = IncrUsn(userId)
			AddNoteUsn(notelist[i].Path, userId, Usn)
		}
	}
	if err := noteList.DeleteCategory(); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	if err := cate.Delete(); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	} else {
		Usn = IncrUsn(userId)
		AddCategoryUsn(userId, cate.Path, Usn)
	}
	tm := time.Now().Unix()
	t := struct {
		Path string
		Usn  int
		Time int64
	}{path, Usn, tm}
	return utils.NewResponse(0, "", t)
}

//DeleteCategoryAndNote delete one category and this category's notes
func DeleteCategoryAndNote(Ob map[string]interface{}, userId string) *utils.Response {
	var cate *model.Category
	var path string
	var Usn int
	if _, ok := Ob["Path"]; ok {
		path = Ob["Path"].(string)
		cate = model.NewCategory(userId, path)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["UpdateTime"]; ok {
		cate.UpdateTime = utils.GetTime(Ob["UpdateTime"])
	} else {
		cate.UpdateTime = utils.NowSecond()
	}
	if !cate.IsExist() {
		return utils.CATEGORY_NOT_EXIST
	}
	if err := cate.Delete(); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	} else {
		Usn = IncrUsn(userId)
		AddCategoryUsn(userId, cate.Path, Usn)
		noteList := model.NewNoteList(path, userId)
		notelist, err := noteList.GetList()
		if err == nil {
			for i := 0; i < len(notelist); i++ {
				Usn = IncrUsn(userId)
				AddNoteUsn(notelist[i].Path, userId, Usn)
			}
		}
		if err := noteList.Delete(cate.UpdateTime); err != nil {
			beego.Error(err)
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}

	}

	tm := time.Now().Unix()
	t := struct {
		Path string
		Usn  int
		Time int64
	}{path, Usn, tm}
	return utils.NewResponse(0, "", t)
}

//GetCategory get all category message
func GetCategory(isDelete bool, userId string) *utils.Response {
	var categoryList []model.Category
	var err error
	cates := model.NewCategoryList(userId)
	if isDelete {
		categoryList, err = cates.GetList()
	} else {
		categoryList, err = cates.GetListDeleted()
	}
	if err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	cateList := convertCategory(categoryList)
	return utils.NewResponse(0, "", cateList)
}

//WebGetCategory get all category message
func WebGetCategory(userId string) interface{} {
	var (
		categoryList []model.Category
		allcount     int
		deletecount  int
		err          error
	)
	notes := model.NewNoteList("", userId)
	cates := model.NewCategoryList(userId)
	categoryList, err = cates.GetList()
	if err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	if allcount, err = notes.GetNoteListCount(model.STATUS_NORMAL); err != nil {
		beego.Error(err)
	}
	if deletecount, err = notes.GetNoteListCount(model.STATUS_DELETED); err != nil {
		beego.Error(err)
	}
	cateList := convertCategory(categoryList)
	t := struct {
		Allcount    int
		CateList    interface{}
		Deletecount int
	}{allcount, cateList, deletecount}
	return utils.NewResponse(utils.SUCCESS_CODE, "", t)
}

func GetAllnote(userId string) interface{} {
	var noteone utils.Note
	noteList := model.NewNoteList("", userId)
	notes, err := noteList.GetAllList()
	if err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	if len(notes) != 0 {
		noteone = NoteResp(notes[0])
	}
	t := struct {
		NoteList interface{}
		NoteOne  interface{}
	}{convertNote(notes), noteone}
	return utils.NewResponse(0, "", t)

}

//GetRecycleBin get recyclebin's all note and category
func GetRecycleBin(userId string) *utils.Response {
	var recycleBins []utils.NoteDeleteList
	var bsm bson.M
	var noteone utils.Note
	bsm = bson.M{}
	notes := model.NewNoteList("", userId)
	notelist, err1 := notes.GetNoteDeletedList(bsm)
	if err1 != nil {
		beego.Debug(err1)
		return utils.NewResponse(utils.SYSTEM_CODE, err1.Error(), nil)
	}
	for i := 0; i < len(notelist); i++ {
		var recycle utils.NoteDeleteList
		recycle.IsShare = notelist[i].IsShare
		recycle.IsVote = notelist[i].IsVote
		recycle.Path = notelist[i].Path
		recycle.CategoryPath = notelist[i].CategoryPath
		recycle.NoteName = notelist[i].Title
		recycle.UpdateTime = notelist[i].UpdateTime.Unix()
		recycle.Type = 1
		recycle.Summary = notelist[i].Summary
		attach := model.NewAttachmentList(notelist[i].Path)
		attachList, err := attach.GetList()
		if err != nil {
			beego.Error(err)
		}
		for i := 0; i < len(attachList); i++ {
			if attachList[i].Status == model.STATUS_DELETED {
				continue
			}
			recycle.AttachTypeList = append(recycle.AttachTypeList, attachList[i].Type)
		}

		recycleBins = append(recycleBins, recycle)
	}
	if len(notelist) != 0 {
		noteone = NoteResp(notelist[0])
	}
	t := struct {
		NoteList interface{}
		NoteOne  interface{}
	}{recycleBins, noteone}
	return utils.NewResponse(0, "", t)
	//	return utils.NewResponse(0, "", recycleBins)
}

//convertCategory 格式转换
func convertCategory(result []model.Category) *[]utils.Category {
	var cates []utils.Category
	for i := 0; i < len(result); i++ {
		var cate utils.Category
		notelist := model.NewNoteList(result[i].Path, result[i].UID)
		notes, err := notelist.GetList()
		if err == nil && notes != nil {
			cate.Count = len(notes)
		}
		cate.Path = result[i].Path
		cate.Name = result[i].Name
		cate.CreateTime = result[i].CreateTime.Unix()
		cate.UpdateTime = result[i].UpdateTime.Unix()
		cates = append(cates, cate)
	}

	return &cates
}

//GetNoteList get one category all note message
func GetNoteList(Ob map[string]interface{}, userId string, isDelete bool) *utils.Response {
	var path string

	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	return Getnotes(path, userId, isDelete)
}
func Getnotes(path, userId string, isDelete bool) *utils.Response {
	var notes []model.Note
	var err error
	var noteone utils.Note
	noteList := model.NewNoteList(path, userId)
	if isDelete {
		notes, err = noteList.GetDeleted()
	} else {
		notes, err = noteList.GetList()
	}
	if err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	if isDelete {
		var noteAll []utils.NoteDeleteList
		for i := 0; i < len(notes); i++ {
			var note utils.NoteDeleteList
			note.IsShare = notes[i].IsShare
			note.IsVote = notes[i].IsVote
			note.Path = notes[i].Path
			note.CategoryPath = notes[i].CategoryPath
			note.NoteName = notes[i].Title
			note.UpdateTime = notes[i].UpdateTime.Unix()
			note.Type = 1
			note.Summary = notes[i].Summary
			attach := model.NewAttachmentList(notes[i].Path)
			attachList, err := attach.GetList()
			if err != nil {
				beego.Error(err)
			}
			for i := 0; i < len(attachList); i++ {
				if attachList[i].Status == model.STATUS_DELETED {
					continue
				}
				note.AttachTypeList = append(note.AttachTypeList, attachList[i].Type)
			}
			noteAll = append(noteAll, note)
		}
		return utils.NewResponse(0, "", noteAll)
	}
	if len(notes) != 0 {
		noteone = NoteResp(notes[0])
	}
	t := struct {
		NoteList interface{}
		NoteOne  interface{}
	}{convertNote(notes), noteone}
	return utils.NewResponse(0, "", t)
	//	return utils.NewResponse(0, "")
}

//SearchName get one category all note message
func SearchName(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		notelist []model.Note
		name     string
		err      error
	)

	if _, ok := Ob["Name"].(string); ok {
		name = Ob["Name"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}

	notes := model.NewNoteList("", userId)
	notelist, err = notes.SearchNoteList(name)
	beego.Debug(notelist)
	if err != nil {
		beego.Debug(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}

	return utils.NewResponse(0, "", convertNote(notelist))

}

//convertNote 格式转换
func convertNote(notes []model.Note) *[]utils.NoteAllList {

	var noteAll []utils.NoteAllList
	for i := 0; i < len(notes); i++ {
		var note utils.NoteAllList
		note.IsShare = notes[i].IsShare
		note.IsVote = notes[i].IsVote
		note.IsPassworded = notes[i].IsPassword
		note.Path = notes[i].Path
		note.CategoryPath = notes[i].CategoryPath

		note.CheckSum = notes[i].CheckSum
		note.NoteName = notes[i].Title
		note.UpdateTime = notes[i].UpdateTime.Unix()
		note.Summary = notes[i].Summary
		attach := model.NewAttachmentList(notes[i].Path)
		attachList, err := attach.GetList()
		if err != nil {
			beego.Error(err)
		}
		for i := 0; i < len(attachList); i++ {
			if attachList[i].Type != "1" && attachList[i].Type != "2" && attachList[i].Type != "3" {
				continue
			}
			if utils.CheckString(note.AttachTypeList, attachList[i].Type) {
				continue
			}
			note.AttachTypeList = append(note.AttachTypeList, attachList[i].Type)
		}
		note.IsPinned = notes[i].IsPinned
		noteAll = append(noteAll, note)
	}
	return &noteAll
}

//UpdateCategory Update one category
func UpdateCategory(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		cate       *model.Category
		path, name string
		Usn        int
	)
	filed := []string{"Name", "UpdateTime", "Status"}
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
		cate = model.NewCategory(userId, path)
		cate.Status = model.STATUS_NORMAL
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["UpdateTime"]; ok {
		cate.UpdateTime = utils.GetTime(Ob["UpdateTime"])
	} else {
		cate.UpdateTime = utils.NowSecond()
	}
	if _, ok := Ob["Name"].(string); ok {
		name = Ob["Name"].(string)
		cate.Name = name
		if cate.Name == "" {
			return utils.CATEGORY_NAME_NOT_NULL
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	//	category, err := cate.GetByNameAndUID(cate.Name, cate.UID)
	//	if category != nil || err != nil {
	//		return utils.CATEGORY_NAME_EXIST
	//	}
	//	if !cate.IsExist() {
	//		return utils.CATEGORY_NOT_EXIST
	//	}
	catetest := model.NewCategory(userId, path)
	group, _ := catetest.GetAllStatus()
	if group.Name != cate.Name {
		category, err := cate.GetByNameAndUID(cate.Name, cate.UID)
		if category != nil || err != nil {
			return utils.CATEGORY_NAME_EXIST
		}
		cates, _ := cate.GetRecycleByNameAndUID(cate.Name, cate.UID)
		if cates != nil {
			//			noteList := model.NewNoteList(cates.Path, userId)
			//			if err := noteList.DeleteCategory(cate.UpdateTime); err != nil {
			//				beego.Error(err)
			//			}
			cates.Remove()
		}
	}
	if err := cate.Update(filed); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	} else {
		Usn = IncrUsn(userId)
		AddCategoryUsn(userId, cate.Path, Usn)
		noteList := model.NewNoteList(path, userId)
		notelist, err := noteList.GetList()
		if err == nil {
			for i := 0; i < len(notelist); i++ {
				Usn = IncrUsn(userId)
				AddNoteUsn(notelist[i].Path, userId, Usn)
			}
			if err := noteList.UpdateCategory(name); err != nil {
				beego.Error(err)
				return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
			}
		} else {
			beego.Error(err)
		}
	}

	tm := time.Now().Unix()
	t := struct {
		Path string
		Usn  int
		Time int64
	}{path, Usn, tm}
	return utils.NewResponse(0, "", t)
}

//AddNoteUsn add one note usn
func AddCategoryUsn(userId, path string, Usn int) {
	cate := model.NewCategory(userId, path)
	filed := []string{"Usn"}
	cate.Usn = Usn
	beego.Debug("inc Usn")
	if err := cate.Update(filed); err != nil {
		beego.Error(err)
	}
	return
}

//RecoverCategory recover one category from recycle bin
func RecoverCategory(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		cate *model.Category
		path string
		Usn  int
	)

	filed := []string{"Status", "UpdateTime"}

	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
		cate = model.NewCategory(userId, path)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	category, err := cate.GetAllStatus()
	if err != nil || category == nil {
		return utils.CATEGORY_NOT_EXIST
	}
	cate.UpdateTime = utils.NowSecond()
	if _, ok := Ob["Name"].(string); ok {
		filed = append(filed, "Name")
		cate.Name = Ob["Name"].(string)
	}
	category, err = cate.GetByNameAndUID(cate.Name, userId)
	if category != nil || err != nil {
		return utils.CATEGORY_NAME_EXIST
	}
	cate.Status = model.STATUS_NORMAL
	if err = cate.Update(filed); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	} else {
		Usn = IncrUsn(userId)
		AddCategoryUsn(userId, cate.Path, Usn)
		noteList := model.NewNoteList(path, userId)
		notelist, err := noteList.GetList()
		if err == nil {
			for i := 0; i < len(notelist); i++ {
				Usn = IncrUsn(userId)
				AddNoteUsn(notelist[i].Path, userId, Usn)
			}
		}
		if err := noteList.Recover(); err != nil {
			beego.Error(err)
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}
	}
	tm := time.Now().Unix()
	t := struct {
		Path string
		Usn  int
		Time int64
	}{path, Usn, tm}
	return utils.NewResponse(0, "", t)

}

//RemoveCategory remove one category from recycle bin
func RemoveCategory(Ob map[string]interface{}, userId string) *utils.Response {
	var cate *model.Category
	var path string
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
		cate = model.NewCategory(userId, path)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	category, err := cate.GetAllStatus()
	if err != nil || category == nil {
		return utils.CATEGORY_NOT_EXIST
	}
	if err = cate.TrashClear(); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	} else {

		noteList := model.NewNoteList(path, userId)
		if err := noteList.TrashClear(); err != nil {
			beego.Error(err)
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}
		notes, err1 := noteList.GetDeleted()
		if notes != nil && err1 == nil {
			for i := 0; i < len(notes); i++ {
				DeleteAttachment(notes[i].Path)
			}
		}
	}
	t := struct {
		Path string
	}{path}
	return utils.NewResponse(0, "", t)

}
