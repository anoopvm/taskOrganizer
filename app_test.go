package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func executeRequest(request *http.Request) *httptest.ResponseRecorder {
	requestResponse := httptest.NewRecorder()
	a.Router.ServeHTTP(requestResponse, request)
	return requestResponse
}

func checkResponseCode(t *testing.T, expectedCode, actualCode int) {
	if expectedCode != actualCode {
		t.Errorf("Wrong code. Expected %d, Actual %d", expectedCode, actualCode)
	}
}

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize("avm", "testing123", "task_organizer")

	code := m.Run()
	os.Exit(code)
}

func TestCreateTask(t *testing.T) {
	payload := []byte(`{"task":"test task by unit test"}`)
	request, err := http.NewRequest("POST", "/task", bytes.NewBuffer(payload))
	if err != nil {
		t.Error("fail")
	}
	RequestResponse := executeRequest(request)
	checkResponseCode(t, http.StatusOK, RequestResponse.Code)
	//t := NewTask()
	//json.Unmarshal(RequestResponse.Body.Bytes(), &t)
}

func TestGetTask(t *testing.T) {

}

func TestListTasks(t *testing.T) {
	request, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Error("fail")
	}
	requestResponse := executeRequest(request)
	checkResponseCode(t, http.StatusOK, requestResponse.Code)
}
