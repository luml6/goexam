package controllers

import (
	"encoding/json"

	"notepad-api/services"

	"github.com/astaxie/beego"
)

type CloudBoxController struct {
	BaseController
}

func getArgs(c *CloudBoxController) map[string]interface{} {

	var args map[string]interface{} = make(map[string]interface{})
	beego.Debug(string(c.Ctx.Input.RequestBody))
	json.Unmarshal(c.Ctx.Input.RequestBody, &args)

	if aid := c.Ctx.Input.Header("AccountID"); aid != "" {
		args["AccountID"] = aid
	}

	return args
}

func (c *CloudBoxController) GetUploadToken() {
	args := getArgs(c)
	cloudServ := services.GetCloudBoxService()
	result := cloudServ.GetUploadToken(args)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *CloudBoxController) GetUploadTokenWithCB() {
	args := getArgs(c)
	cloudServ := services.GetCloudBoxService()
	result := cloudServ.GetUploadTokenWithCallback(args)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *CloudBoxController) GetDownloadURL() {
	args := getArgs(c)
	cloudServ := services.GetCloudBoxService()
	result := cloudServ.GetDownloadURL(args)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *CloudBoxController) CallBack() {

	var Ob map[string]interface{}
	beego.Debug(string(c.Ctx.Input.RequestBody))
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	beego.Debug(Ob)
	result := services.AttachCallback(Ob)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
