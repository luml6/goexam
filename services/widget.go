package services

import (
	"notepad-api/model"
	"notepad-api/utils"
)

//GetWidgetList get all widget
func GetWidgetList() *utils.Response {
	widget := new(model.Widget)
	widgetList, err := widget.GetList()
	if err == nil || widgetList != nil {
		return utils.NewResponse(0, "", widgetList)
	}
	return utils.NewResponse(0, "", nil)
}

//AddWidget add one widget
func AddWidget(Ob map[string]interface{}) *utils.Response {
	var widget model.Widget
	widget.ID = utils.CreateId()
	if _, ok := Ob["Name"].(string); ok {
		widget.Name = Ob["Name"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["CateName"].(string); ok {
		widget.CateName = Ob["CateName"].(string)
	}
	if _, ok := Ob["Price"].(float64); ok {
		widget.Price = int(Ob["Price"].(float64))
	}

	if err := widget.Add(); err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		ID string
	}{widget.ID}
	return utils.NewResponse(0, "", t)
}

//UpdateWidget update one widget
func UpdateWidget(Ob map[string]interface{}) *utils.Response {
	var widget *model.Widget
	var filed []string
	if _, ok := Ob["ID"].(string); ok {
		ID := Ob["ID"].(string)
		widget = model.NewWidget(ID)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if !widget.IsExist() {
		return utils.WIDGET_NOT_EXIST
	}
	if _, ok := Ob["Name"].(string); ok {
		widget.Name = Ob["Name"].(string)
		filed = append(filed, "Name")
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["CateName"].(string); ok {
		widget.CateName = Ob["CateName"].(string)
		filed = append(filed, "CateName")
	}
	if _, ok := Ob["Price"].(float64); ok {
		widget.Price = int(Ob["Price"].(float64))
		filed = append(filed, "Price")
	}
	widget.UpdateTime = utils.NowSecond()
	filed = append(filed, "UpdateTime")
	if err := widget.Update(filed); err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		ID string
	}{widget.ID}
	return utils.NewResponse(0, "", t)
}

//DeleteWidget delete one widget
func DeleteWidget(Ob map[string]interface{}) *utils.Response {
	var ID string
	var widget *model.Widget
	if _, ok := Ob["ID"].(string); ok {
		ID = Ob["ID"].(string)
		widget = model.NewWidget(ID)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if !widget.IsExist() {
		return utils.WIDGET_NOT_EXIST
	}
	if err := widget.Delete(); err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		ID string
	}{ID}
	return utils.NewResponse(0, "", t)
}
