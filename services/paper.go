package services

import (
	"notepad-api/model"
	"notepad-api/utils"
)

//GetPaperList get all paper
func GetPaperList() *utils.Response {
	paper := new(model.Background)
	paperList, err := paper.GetList()
	if err == nil || paperList != nil {
		return utils.NewResponse(0, "", paperList)
	}
	return utils.NewResponse(0, "", nil)
}

//AddPaper add one paper
func AddPaper(Ob map[string]interface{}) *utils.Response {
	var paper model.Background
	paper.ID = utils.CreateId()
	if _, ok := Ob["Name"].(string); ok {
		paper.Name = Ob["Name"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["CateName"].(string); ok {
		paper.CateName = Ob["CateName"].(string)
	}
	if _, ok := Ob["Content"].(string); ok {
		paper.Content = Ob["Content"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if err := paper.Add(); err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		ID string
	}{paper.ID}
	return utils.NewResponse(0, "", t)
}

//UpdatePaper update one paper
func UpdatePaper(Ob map[string]interface{}) *utils.Response {
	var paper *model.Background
	var filed []string
	if _, ok := Ob["ID"].(string); ok {
		ID := Ob["ID"].(string)
		paper = model.NewBackground(ID)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if !paper.IsExist() {
		return utils.PAPER_NOT_EXIST
	}
	if _, ok := Ob["Name"].(string); ok {
		paper.Name = Ob["Name"].(string)
		filed = append(filed, "Name")
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["CateName"].(string); ok {
		paper.CateName = Ob["CateName"].(string)
		filed = append(filed, "CateName")
	}
	if _, ok := Ob["Content"].(string); ok {
		paper.Content = Ob["Content"].(string)
		filed = append(filed, "Content")
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	paper.UpdateTime = utils.NowSecond()
	filed = append(filed, "UpdateTime")
	if err := paper.Update(filed); err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		ID string
	}{paper.ID}
	return utils.NewResponse(0, "", t)
}

//DeletePaper update one paper
func DeletePaper(Ob map[string]interface{}) *utils.Response {
	var paper *model.Background
	var ID string
	if _, ok := Ob["ID"].(string); ok {
		ID = Ob["ID"].(string)
		paper = model.NewBackground(ID)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if !paper.IsExist() {
		return utils.PAPER_NOT_EXIST
	}
	if err := paper.Delete(); err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		ID string
	}{ID}
	return utils.NewResponse(0, "", t)
}
