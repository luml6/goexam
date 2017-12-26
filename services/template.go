package services

import (
	"notepad-api/model"
	"notepad-api/utils"
)

//GetTemplateList get all template
func GetTemplateList() *utils.Response {
	template := new(model.Template)
	temList, err := template.GetList()
	if err == nil || temList != nil {
		return utils.NewResponse(0, "", temList)
	}
	return utils.NewResponse(0, "", nil)
}

//AddTemplate add one template
func AddTemplate(Ob map[string]interface{}) *utils.Response {
	var template model.Template
	template.ID = utils.CreateId()
	if _, ok := Ob["Name"].(string); ok {
		template.Name = Ob["Name"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["CateName"].(string); ok {
		template.CateName = Ob["CateName"].(string)
	}
	if _, ok := Ob["Content"].(string); ok {
		template.Content = Ob["Content"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["Price"].(float64); ok {
		template.Price = int(Ob["Price"].(float64))
	}
	if err := template.Add(); err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		Path string
	}{template.ID}
	return utils.NewResponse(0, "", t)
}

//UpdateTemplate update one template
func UpdateTemplate(Ob map[string]interface{}) *utils.Response {
	var template *model.Template
	var filed []string
	if _, ok := Ob["ID"].(string); ok {
		ID := Ob["ID"].(string)
		template = model.NewTemplate(ID)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if !template.IsExist() {
		return utils.TEMPLATE_NOT_EXIST
	}
	if _, ok := Ob["Name"].(string); ok {
		template.Name = Ob["Name"].(string)
		filed = append(filed, "Name")
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["CateName"].(string); ok {
		template.CateName = Ob["CateName"].(string)
		filed = append(filed, "CateName")
	}
	if _, ok := Ob["Content"].(string); ok {
		template.Content = Ob["Content"].(string)
		filed = append(filed, "Content")
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["Price"].(float64); ok {
		template.Price = int(Ob["Price"].(float64))
		filed = append(filed, "Price")
	}
	template.UpdateTime = utils.NowSecond()
	filed = append(filed, "UpdateTime")
	if err := template.Update(filed); err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		Path string
	}{template.ID}
	return utils.NewResponse(0, "", t)
}

//DeleteTemplate update one paper
func DeleteTemplate(Ob map[string]interface{}) *utils.Response {
	var tem *model.Template
	var ID string
	if _, ok := Ob["ID"].(string); ok {
		ID = Ob["ID"].(string)
		tem = model.NewTemplate(ID)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if !tem.IsExist() {
		return utils.TEMPLATE_NOT_EXIST
	}
	if err := tem.Delete(); err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		ID string
	}{ID}
	return utils.NewResponse(0, "", t)
}
