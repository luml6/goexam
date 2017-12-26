package tests

import (
	"notepad-api/services"
	"testing"
)

func TestCreateCategory(t *testing.T) {
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	ob := make(map[string]interface{})
	ob["Name"] = "test"
	ob["CreateTime"] = 1474617746
	resp := services.CreateCategory(ob, UserId)
	t.Log(resp)

}

func TestDeleteCategory(t *testing.T) {
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	ob := make(map[string]interface{})
	ob["Path"] = "8d887fa9-41a8-432b-bb29-7edfcf45bdff"
	ob["UpdateTime"] = 1474617746
	resp := services.DeleteCategory(ob, UserId)
	t.Log(resp)

}

func TestDeleteCategoryAndNote(t *testing.T) {
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	ob := make(map[string]interface{})
	ob["Path"] = "8d887fa9-41a8-432b-bb29-7edfcf45bdff"
	ob["UpdateTime"] = 1474617746
	resp := services.DeleteCategory(ob, UserId)
	t.Log(resp)
}

func TestGetCategory(t *testing.T) {
	isDelete := true
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	resp := services.GetCategory(isDelete, UserId)
	t.Log(resp)
}

func TestGetNoteList(t *testing.T) {
	UserId := "587ce63a-0fde-49ec-a037-1fd3193206f3"
	ob := make(map[string]interface{})
	ob["Path"] = "8d887fa9-41a8-432b-bb29-7edfcf45bdff"
	resp := services.GetNoteList(ob, UserId, true)
	t.Log(resp)
}

func TestGetOtherNoteList(t *testing.T) {
	UserId := "587ce63a-0fde-49ec-a037-1fd3193206f3"
	resp := services.GetOtherNoteList(UserId)
	t.Log(resp)
}

func TestUpdateCategory(t *testing.T) {
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	ob := make(map[string]interface{})
	ob["Path"] = "8d887fa9-41a8-432b-bb29-7edfcf45bdff"
	ob["UpdateTime"] = 1474617746
	ob["Name"] = "test"
	resp := services.UpdateCategory(ob, UserId)
	t.Log(resp)
}
