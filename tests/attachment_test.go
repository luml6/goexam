package tests

import (
	"testing"
)

// /V1.0/Attachment/Add
func TestAdd(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Attachment/Add"
	method := "POST"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	params := `{
				  "FileType": "string",
				  "NoteId": "string",
				  "Uuid": "string",
				  "ObjectKey": "string"
				}`
	resp, err := HttpRequest(url, method, params, UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}
