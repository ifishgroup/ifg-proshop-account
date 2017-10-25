// +build integration

package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ifishgroup/ifg-proshop-account/db"
)

var conn *sql.DB
var err error
var host = "postgres"
var user = "postgres"
var pass = "postgres"

func TestIntegrationHandleGet(t *testing.T) {
	mux := http.NewServeMux()
	account := &db.Account{Db: newDbConnection()}
	mux.HandleFunc("/account/", handleRequest(account))
	writer := httptest.NewRecorder()

	postAccountData(t, mux, writer)

	request, _ := http.NewRequest("GET", fmt.Sprintf("/account/%d", account.Id), nil)

	mux.ServeHTTP(writer, request)
	if writer.Code != 200 {
		t.Errorf("GET response code is %v", writer.Code)
	}

	json.Unmarshal(writer.Body.Bytes(), &account)

	if account.Name != "foo" {
		t.Error("Cannot retrieve JSON account")
	}

	conn.Close()
}

func TestIntgerationHandlePut(t *testing.T) {
	mux := http.NewServeMux()
	account := &db.Account{Db: newDbConnection()}
	mux.HandleFunc("/account/", handleRequest(account))
	writer := httptest.NewRecorder()

	postAccountData(t, mux, writer)

	jsonUpdate := strings.NewReader(`{"name":"foo","email":"foo@gmail.com"}`)
	request, _ := http.NewRequest("PUT", fmt.Sprintf("/account/%d", account.Id), jsonUpdate)
	mux.ServeHTTP(writer, request)

	if account.Email != "foo@gmail.com" {
		t.Error("Email is not correct", account.Name)
	}
	conn.Close()
}

func TestIntegrationHandlePost(t *testing.T) {
	mux := http.NewServeMux()
	account := &db.Account{Db: newDbConnection()}
	mux.HandleFunc("/account/", handleRequest(account))
	writer := httptest.NewRecorder()

	json := strings.NewReader(`{"name":"foo","email":"foo@bar.com"}`)
	request, _ := http.NewRequest("POST", "/account/", json)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	cleanup()
	conn.Close()
}

func TestIntegrationHandleDelete(t *testing.T) {
	mux := http.NewServeMux()
	account := &db.Account{Db: newDbConnection()}
	mux.HandleFunc("/account/", handleRequest(account))
	writer := httptest.NewRecorder()

	postAccountData(t, mux, writer)

	request, _ := http.NewRequest("DELETE", fmt.Sprintf("/account/%d", account.Id), nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	conn.Close()
}

func postAccountData(t *testing.T, mux *http.ServeMux, writer *httptest.ResponseRecorder) {
	json := strings.NewReader(`{"name":"foo","email":"foo@bar.com"}`)
	request, _ := http.NewRequest("POST", "/account/", json)
	mux.ServeHTTP(writer, request)
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func newDbConnection() *sql.DB {
	conn, err = db.NewDbConnectionWithParams(host, user, pass)
	if err != nil {
		panic(err)
	}
	return conn
}

func cleanup() {
	conn.Exec("delete from accounts")
}
