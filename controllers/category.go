package controllers

import (
	"notepad-api/services"

	"encoding/json"

	"github.com/astaxie/beego"
)

type CategoryController struct {
	BaseController
}

//All get all category message
func (c *CategoryController) All() {
	AccountID := c.Ctx.Input.Header("AccountID")
	beego.Debug(AccountID)
	cates := services.GetCategory(true, AccountID)
	beego.Debug(cates)
	c.Data["json"] = cates
	c.ServeJSON()
}

//WebAll get all category message
func (c *CategoryController) WebAll() {
	//	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	beego.Debug(AccountID)

	cates := services.WebGetCategory(AccountID)
	beego.Debug(cates)
	c.Data["json"] = cates
	c.ServeJSON()
}

func (c *CategoryController) Search() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.SearchName(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

//NoteList get one category's all note message
func (c *CategoryController) NoteList() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.GetNoteList(Ob, AccountID, false)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

//NoteDeletedList get one deleted category's all note message
func (c *CategoryController) NoteDeletedList() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.GetNoteList(Ob, AccountID, true)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

//All get all category message
func (c *CategoryController) AllDelete() {
	AccountID := c.Ctx.Input.Header("AccountID")
	cates := services.GetRecycleBin(AccountID)
	beego.Debug(cates)

	c.Data["json"] = cates
	c.ServeJSON()
}

//Create create a new category
func (c *CategoryController) Create() {

	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.CreateCategory(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

//Update update one category message
func (c *CategoryController) Update() {

	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.UpdateCategory(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

//Delete delete one category
func (c *CategoryController) Delete() {

	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.DeleteCategory(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

//DeleteAll delete one category and delete this category's notes
func (c *CategoryController) DeleteAll() {

	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.DeleteCategoryAndNote(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

func (c *CategoryController) RecoverCategory() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.RecoverCategory(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *CategoryController) RemoveCategory() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.RemoveCategory(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
