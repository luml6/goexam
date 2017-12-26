package controllers

import (
	"notepad-api/services"
	//	"notepad-api/utils"

	"encoding/json"

	"github.com/astaxie/beego"
)

type TemplateController struct {
	beego.Controller
}

func (c *TemplateController) AddTemplate() {

	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.AddTemplate(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *TemplateController) TemplateList() {

	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.GetTemplateList()
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *TemplateController) UpdateTemplate() {
	var Ob map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	beego.Debug(err)
	result := services.UpdateTemplate(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *TemplateController) DeleteTemplate() {
	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.DeleteTemplate(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
