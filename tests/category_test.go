package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func HttpRequest(url, method, params, userId string) ([]byte, error) {
	var req *http.Request
	var err error
	client := &http.Client{}
	if params != "" {
		req, err = http.NewRequest(method, url, strings.NewReader(params))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		fmt.Println(err)
		// handle error
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AccountID", userId)

	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

// /V1.0/Category/All
func TestAll(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Category/All"
	method := "GET"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	resp, err := HttpRequest(url, method, " ", UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}

// /V1.0/Category/NoteList
func TestNoteList(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Category/NoteList"
	method := "POST"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	params := `{"Path": "8d887fa9-41a8-432b-bb29-7edfcf45bdff"}`
	//	params, _ := json.Marshal(paramJson)
	resp, err := HttpRequest(url, method, params, UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}

// /V1.0/Category/Create
func TestCreateCategory(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Category/Create"
	method := "POST"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	params := `{"Name": "categoryTest1","CreateTime": "2016-01-02 00:00:00"}`
	resp, err := HttpRequest(url, method, params, UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}

// /V1.0/Category/Update
func TestUpdateCategory(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Category/Update"
	method := "POST"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	params := `{"Path":"7e2bba8d-257f-4325-b3b6-43775cc6f4c8","Name": "categoryTest1","UpdateTime": "2016-01-02 00:00:00"}`
	resp, err := HttpRequest(url, method, params, UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}

// /V1.0/Category/Delete
func TestDeleteCategory(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Category/Delete"
	method := "POST"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	params := `{"Path":"7e2bba8d-257f-4325-b3b6-43775cc6f4c8","UpdateTime": "2016-01-02 00:00:00"}`
	resp, err := HttpRequest(url, method, params, UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}

// /V1.0/Category/DeleteAll
func TestDeleteAll(t *testing.T) {
	url := "http://10.0.50.185:8080/V1.0/Category/DeleteAll"
	method := "POST"
	UserId := "19b485ca-1c32-4df6-b9be-b24ca1f2f83d"
	params := `{"Path":"7e2bba8d-257f-4325-b3b6-43775cc6f4c8","UpdateTime": "2016-01-02 00:00:00"}`
	resp, err := HttpRequest(url, method, params, UserId)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(resp))
}
