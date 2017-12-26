package controllers

import (
	"encoding/json"
	"notepad-api/conf"
	"notepad-api/services"

	"github.com/astaxie/beego"
)

type UserController struct {
	BaseController
}

func (c *UserController) GetToken() {
	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	t := services.PushAccount(Ob)
	beego.Debug(t)
	c.Data["json"] = t
	c.ServeJSON()
}

func (c *UserController) AppGetToken() {
	var Ob map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	//	AccountID := c.Ctx.Input.Header("AccountID")
	//	RegisterID := c.Ctx.Input.Header("RegisterID")
	t := services.AppGetToken(Ob)
	beego.Debug(t)
	c.Data["json"] = t
	c.ServeJSON()
}
func (c *UserController) AllTrash() {
	AccountID := c.Ctx.Input.Header("AccountID")
	t := services.AllTrash(AccountID)
	beego.Debug(t)
	c.Data["json"] = t
	c.ServeJSON()
}

func (c *UserController) GetUserInfo() {
	AccountID := c.Ctx.Input.Header("AccountID")
	t := services.GetUserInfo(AccountID)
	beego.Debug(t)
	c.Data["json"] = t
	c.ServeJSON()
}
func (c *UserController) LoginOut() {
	AccountID := c.Ctx.Input.Header("AccountID")
	token := c.Ctx.Input.Header("token")
	t := services.LoginOut(AccountID, token)
	if t {
		c.Redirect(conf.Conf.LoginOutUrl, 302)
		return
	}
	re := struct {
		IsSuccess bool
	}{false}
	c.Data["json"] = re
	c.ServeJSON()
}

func (c *UserController) Cancellation() {
	AccountID := c.Ctx.Input.Header("AccountID")
	RegisterID := c.Ctx.Input.Header("PushID")
	t := services.Cancellation(AccountID, RegisterID)
	re := struct {
		IsSuccess bool
	}{t}
	c.Data["json"] = re
	c.ServeJSON()
}
