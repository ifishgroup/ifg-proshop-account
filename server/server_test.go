package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ifishgroup/ifg-proshop-account/db"
)

var mux *http.ServeMux
var writer *httptest.ResponseRecorder
var account *db.FakeAccount

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	mux = http.NewServeMux()
	account = &db.FakeAccount{}
	mux.HandleFunc("/account/", handleRequest(account))
	writer = httptest.NewRecorder()
}

func tearDown() {
}

func TestHandleGet(t *testing.T) {
	request, _ := http.NewRequest("GET", "/account/1", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	var account db.Account
	json.Unmarshal(writer.Body.Bytes(), &account)

	if account.Id != 1 {
		t.Error("Cannot retrieve JSON account")
	}
}

func TestHandlePut(t *testing.T) {
	json := strings.NewReader(`{"name":"foo","email":"foo@bar.com"}`)
	request, _ := http.NewRequest("PUT", "/account/1", json)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	if account.Name != "foo" {
		t.Error("Name is not correct", account.Name)
	}
}

func TestHandlePost(t *testing.T) {
	json := strings.NewReader(`{"name":"foo","email":"foo@bar.com"}`)
	request, _ := http.NewRequest("POST", "/account/1", json)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestHandleDelete(t *testing.T) {
	request, _ := http.NewRequest("DELETE", "/account/1", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func printJson(account db.Account) {
	prettyPost, err := json.MarshalIndent(&account, "", " ")
	if err != nil {
		fmt.Println("Error pretty parsing json")
		return
	}
	fmt.Println(string(prettyPost))
}
