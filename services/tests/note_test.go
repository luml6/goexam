package tests

import (
	"notepad-api/services"
	"testing"
)

func TestCreate(t *testing.T) {
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	ob := make(map[string]interface{})
	ob["CategoryPath"] = "8d887fa9-41a8-432b-bb29-7edfcf45bdff"
	ob["Title"] = "test"
	ob["Content"] = "test"
	ob["Author"] = "test"
	ob["Summary"] = "test"
	ob["CheckSum"] = "test"
	ob["CreateTime"] = 1474617746
	resp := services.Create(ob, UserId)
	t.Log(resp)

}

func TestDelete(t *testing.T) {
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	ob := make(map[string]interface{})
	ob["Path"] = "4a99cfcf-aed7-440b-9a25-57fe802d6420"
	ob["UpdateTime"] = 1474617746
	resp := services.Delete(ob, UserId)
	t.Log(resp)

}

func TestDeleteAttach(t *testing.T) {
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	ob := make(map[string]interface{})
	ob["Uuid"] = "123456"
	resp := services.DeleteAttach(ob, UserId)
	t.Log(resp)
}

func TestGetNote(t *testing.T) {
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	ob := make(map[string]interface{})
	ob["Path"] = "4a99cfcf-aed7-440b-9a25-57fe802d6420"
	resp := services.GetNote(ob, UserId, 0)
	t.Log(resp)
}

func TestMove(t *testing.T) {
	UserId := "587ce63a-0fde-49ec-a037-1fd3193206f3"
	ob := make(map[string]interface{})
	ob["Path"] = "4a99cfcf-aed7-440b-9a25-57fe802d6420"
	ob["CategoryPath"] = "1c91fac9-c758-4304-8c56-4e022a254aaa"
	resp := services.Move(ob, UserId)
	t.Log(resp)
}

func TestUpdate(t *testing.T) {
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	ob := make(map[string]interface{})
	ob["Path"] = "4a99cfcf-aed7-440b-9a25-57fe802d6420"
	ob["Title"] = "test"
	ob["Content"] = "test"
	ob["Author"] = "test"
	ob["Summary"] = "test"
	ob["CheckSum"] = "test"
	ob["UpdateTime"] = 1474617746
	resp := services.Update(ob, UserId)
	t.Log(resp)
}
