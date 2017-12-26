package controllers

import (
	"notepad-api/services"

	"encoding/json"

	"github.com/astaxie/beego"
)

type NoteController struct {
	BaseController
}

// Create new note
func (c *NoteController) Create() {

	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.Create(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

func (c *NoteController) UpdateByV2() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.UpdateByV2(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

//Update note message
func (c *NoteController) Update() {

	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.Update(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

//GetNote get note content
func (c *NoteController) GetNote() {

	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")

	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.GetNote(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *NoteController) GetNoteDeleted() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")

	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.GetNote(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

//GetNote get note content
func (c *NoteController) GetNoteID() {
	AccountID := c.Ctx.Input.Header("AccountID")
	result := services.GetNoteID(AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}

//Delete note
func (c *NoteController) Delete() {

	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.Delete(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

//DeleteByV2 note
func (c *NoteController) DeleteByV2() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.DeleteByV2(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

//Move note to next category
func (c *NoteController) Move() {

	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.Move(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

//NoCategory get all note message without category
func (c *NoteController) NoCategory() {

	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")

	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.GetOtherNoteList(AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

//All get all note message with one AccountID
func (c *NoteController) AllNote() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.GetAllNoteList(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

//DeleteAttach delete one note attach
func (c *NoteController) DeleteAttach() {

	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.DeleteAttach(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()

}

func (c *NoteController) RecoverNote() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.RecoverNote(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *NoteController) RemoveNote() {
	var Ob map[string]interface{}
	AccountID := c.Ctx.Input.Header("AccountID")
	json.Unmarshal(c.Ctx.Input.RequestBody, &Ob)
	result := services.RemoveNote(Ob, AccountID)
	beego.Debug(result)
	c.Data["json"] = result
	c.ServeJSON()
}
