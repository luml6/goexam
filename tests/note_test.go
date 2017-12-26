package tests

import (
	"testing"
)

// /V1.0/Note/NoCategory
func TestNoCategory(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Note/NoCategory"
	method := "GET"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	resp, err := HttpRequest(url, method, " ", UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}

// /V1.0/Note/Get
func TestNoteGet(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Note/Get"
	method := "POST"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	params := `{"Path": "b156f55f-0645-4473-8279-e476c1dd15c3"}`
	resp, err := HttpRequest(url, method, params, UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}

// /V1.0/Note/Create
func TestCreateNote(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Note/Create"
	method := "POST"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	params := `{
				  "CategoryPath": "2c41385f-275c-4da1-9956-43f2d9753591",
				  "Author": "test",
				  "Title": "test",
				  "Summary": "test",
				  "Content": "test",
				  "CheckSum": "test",
				  "CreateTime": 0
				}`
	resp, err := HttpRequest(url, method, params, UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}

// /V1.0/Note/Update
func TestUpdateNote(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Note/Update"
	method := "POST"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	params := `{
				  "Path": "b156f55f-0645-4473-8279-e476c1dd15c3",
				  "Author": "string",
				  "Title": "string",
				  "Summary": "string",
				  "CheckSum": "string",
				  "Content": "string",
				  "UpdateTime": 0
				}`
	resp, err := HttpRequest(url, method, params, UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}

// /V1.0/Note/Delete
func TestDeleteNote(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Note/Delete"
	method := "POST"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	params := `{"Path":"b156f55f-0645-4473-8279-e476c1dd15c3","UpdateTime": "2016-01-02 00:00:00"}`
	resp, err := HttpRequest(url, method, params, UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}

// /V1.0/Note/DeleteAttach
func TestDeleteAttach(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Note/DeleteAttach"
	method := "POST"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	params := `{"Uuid":"7e2bba8d-257f-4325-b3b6-43775cc6f4c8"}`
	resp, err := HttpRequest(url, method, params, UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}

// /V1.0/Note/Move
func TestMove(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Note/Move"
	method := "POST"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	params := `{"CategoryPath":"7e2bba8d-257f-4325-b3b6-43775cc6f4c8","Path": "b156f55f-0645-4473-8279-e476c1dd15c3"}`
	resp, err := HttpRequest(url, method, params, UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}
