package controllers

import (
	//	"bytes"
	"encoding/json"
	//	"html/template"
	"notepad-api/services"

	"github.com/astaxie/beego"
)

type ShareController struct {
	BaseController
}

func (c *ShareController) Publish() {
	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	AccountID := c.Ctx.Input.Header("AccountID")
	result := services.AddShare(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *ShareController) Show() {
	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.GetShareNote(Ob)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *ShareController) GetShareMessage() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.GetShareMessage(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *ShareController) GetQuestion() {
	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.GetQuestion(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *ShareController) GetShare() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.GetShareUrl(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *ShareController) CancleShare() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.CancleShare(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
