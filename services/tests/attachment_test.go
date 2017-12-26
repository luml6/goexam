package tests

import (
	"notepad-api/services"
	"testing"
)

func TestUpload(t *testing.T) {
	noteId := "0b017253-f93f-4e9c-b91c-8f8377de2318"
	uuid := "123456"
	fileType := "img"
	fileName := "test.jpg"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	resp := services.Upload(noteId, uuid, fileType, fileName, UserId)
	t.Log(resp)
}

func TestAdd(t *testing.T) {
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	ob := make(map[string]interface{})
	ob["Uuid"] = "123456"
	ob["FileType"] = "img"
	ob["FileName"] = "test.jpg"
	ob["NoteId"] = "0b017253-f93f-4e9c-b91c-8f8377de2318"
	ob["ObjectKey"] = "www.baidu.com"
	resp := services.Add(ob, UserId)
	t.Log(resp)

}
