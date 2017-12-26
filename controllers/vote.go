package controllers

import (
	//	"bytes"
	"encoding/json"
	//	"html/template"
	"notepad-api/services"

	"github.com/astaxie/beego"
)

type VoteController struct {
	BaseController
}

func (c *VoteController) VoteAdd() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.VoteAdd(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *VoteController) VoteUpdate() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.VoteUpdate(Ob, AccountID)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *VoteController) VoteGet() {
	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.VoteGet(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *VoteController) VoteGetByNote() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.VoteGetByNote(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *VoteController) CancleVote() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.CancleVote(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *VoteController) ClickVote() {
	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.ClickVote(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
