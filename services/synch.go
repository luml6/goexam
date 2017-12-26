package services

import (
	"notepad-api/model"
	"notepad-api/utils"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//GetUserUsn get one user's Usn
func GetUserUsn(AccountID string) *utils.Response {
	user := model.NewUser(AccountID)
	t := make(map[string]interface{})
	key := utils.USNKEY + AccountID
	if usn, err := utils.Cache.Get(key).Result(); err == nil {
		num, _ := strconv.Atoi(usn)
		t["Usn"] = num
	} else {
		if _, err := user.Get(); err == nil {
			t["Usn"] = user.Usn
		} else {
			t["Usn"] = 0
		}
	}
	t["Time"] = time.Now().Unix()
	return utils.NewResponse(0, "", t)
}

func GetSyncNotes(AccountID string, ob map[string]interface{}) *utils.Response {
	var afterUsn int
	var syncs []utils.SyncResp
	var bsm bson.M
	if usn, ok := ob["Usn"].(float64); ok {
		afterUsn = int(usn)
	}
	notes := model.NewNoteList("", AccountID)
	//	if afterUsn == 0 {
	//		bsm = bson.M{"Status": model.STATUS_NORMAL, "Usn": bson.M{"$gt": afterUsn}}
	//	} else {
	bsm = bson.M{"Usn": bson.M{"$gt": afterUsn}}
	//	}

	notelist, err := notes.GetSyncNoteList(bsm)
	if err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), "")
	}
	for i := 0; i < len(notelist); i++ {
		var syncresp utils.SyncResp
		syncresp.Path = notelist[i].Path
		syncresp.CheckSum = notelist[i].CheckSum
		syncresp.CategoryPath = notelist[i].CategoryPath
		syncresp.Title = notelist[i].Title
		syncresp.Usn = notelist[i].Usn
		syncresp.UpdateTime = notelist[i].UpdateTime.Unix()
		syncresp.IsPinned = notelist[i].IsPinned
		syncresp.IsVote = notelist[i].IsVote
		syncresp.IsPassworded = notelist[i].IsPassword
		if notelist[i].Status != model.STATUS_NORMAL {
			syncresp.IsDeleted = true
		}
		syncs = append(syncs, syncresp)
	}
	return utils.NewResponse(0, "", syncs)

}

func GetSyncCategorys(AccountID string, ob map[string]interface{}) *utils.Response {
	var afterUsn int
	var syncs []utils.SyncResp
	if usn, ok := ob["Usn"].(float64); ok {
		afterUsn = int(usn)
	}
	cates := model.NewCategoryList(AccountID)
	catelist, err := cates.GetSynchList(afterUsn)
	if err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), "")
	}
	for i := 0; i < len(catelist); i++ {
		var syncresp utils.SyncResp
		syncresp.Path = catelist[i].Path
		syncresp.UpdateTime = catelist[i].UpdateTime.Unix()
		syncresp.Name = catelist[i].Name
		if catelist[i].Status != model.STATUS_NORMAL {
			syncresp.IsDeleted = true
		}
		syncs = append(syncs, syncresp)
	}
	return utils.NewResponse(0, "", syncs)
}
