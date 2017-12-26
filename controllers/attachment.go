package controllers

import (
	"notepad-api/services"
	"notepad-api/utils"

	"encoding/json"

	"github.com/astaxie/beego"
)

type AttachmentController struct {
	BaseController
}

//Upload new attachment
func (c *AttachmentController) Upload() {
	var result *utils.Response
	beego.Debug(c.Ctx.Input.Params())
	AccountID := c.Ctx.Input.Header("AccountID")
	fileType := c.Input().Get("fileType")
	noteId := c.Input().Get("noteId")
	uuid := c.Input().Get("uuid")
	f, h, err := c.GetFile("File")
	path := beego.AppConfig.String("uploadPath")
	if err != nil {
		beego.Debug("getfile err ", err)
	} else {
		defer f.Close()
		err = c.SaveToFile("File", path+h.Filename)
		if err == nil {
			result = services.Upload(noteId, uuid, fileType, h.Filename, AccountID)
		} else {
			beego.Debug(err)
		}
	}
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *AttachmentController) Add() {

	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	beego.Debug(string(c.Ctx.Input.RequestBody))
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.Add(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *AttachmentController) AddPublic() {

	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	beego.Debug(string(c.Ctx.Input.RequestBody))
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.AddPublic(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *AttachmentController) DownLoad() {
	fileName := c.Input().Get("id")
	path := beego.AppConfig.String("uploadPath")
	c.Ctx.Output.Download(path + fileName)
}

func (c *AttachmentController) Check() {
	var Ob map[string]interface{}
	beego.Debug(string(c.Ctx.Input.RequestBody))
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.Check(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
