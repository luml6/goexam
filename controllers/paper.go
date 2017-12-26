package controllers

import (
	"notepad-api/services"
	//	"notepad-api/utils"

	"encoding/json"

	"github.com/astaxie/beego"
)

type PaperController struct {
	beego.Controller
}

func (c *PaperController) Paperlist() {

	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.GetPaperList()
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *PaperController) AddPaper() {
	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.AddPaper(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *PaperController) UpdatePaper() {
	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.UpdatePaper(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *PaperController) DeletePaper() {
	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.DeletePaper(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
