package controllers

import (
	"encoding/json"

	"notepad-api/services"

	"github.com/astaxie/beego"
)

type SynchController struct {
	BaseController
}

func (c *SynchController) GetUserUsn() {
	AccountID := c.Ctx.Input.Header("AccountID")
	t := services.GetUserUsn(AccountID)
	beego.Debug(t)
	c.Data["json"] = t
	c.ServeJSON()
}

func (c *SynchController) GetSyncNotes() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	t := services.GetSyncNotes(AccountID, Ob)
	beego.Debug(t)
	c.Data["json"] = t
	c.ServeJSON()
}

func (c *SynchController) GetSyncCategorys() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	t := services.GetSyncCategorys(AccountID, Ob)
	beego.Debug(t)
	c.Data["json"] = t
	c.ServeJSON()
}
