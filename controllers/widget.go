package controllers

import (
	"notepad-api/services"
	//	"notepad-api/utils"

	"encoding/json"

	"github.com/astaxie/beego"
)

type WidgetController struct {
	beego.Controller
}

func (c *WidgetController) Widgetlist() {

	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.GetWidgetList()
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *WidgetController) AddWidget() {
	var Ob map[string]interface{}
	beego.Debug(string(c.Ctx.Input.RequestBody))
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.AddWidget(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *WidgetController) UpdateWidget() {
	var Ob map[string]interface{}
	beego.Debug(string(c.Ctx.Input.RequestBody))
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.UpdateWidget(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *WidgetController) DeleteWidget() {
	var Ob map[string]interface{}
	beego.Debug(string(c.Ctx.Input.RequestBody))
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.DeleteWidget(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
