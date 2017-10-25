package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	"github.com/ifishgroup/ifg-proshop-account/db"
	_ "github.com/lib/pq"
)

func NewServer(connection *sql.DB) *http.Server {
	server := http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("/account/", handleRequest(db.NewAccount(connection)))
	return &server
}

func handleRequest(t db.CrudRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		switch r.Method {
		case "GET":
			err = handleGet(w, r, t)
		case "POST":
			err = handlePost(w, r, t)
		case "PUT":
			err = handlePut(w, r, t)
		case "DELETE":
			err = handleDelete(w, r, t)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleGet(w http.ResponseWriter, r *http.Request, repo db.CrudRepository) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	err = repo.Fetch(id)
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(repo, "", "\t\t")
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

func handlePost(w http.ResponseWriter, r *http.Request, repo db.CrudRepository) (err error) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	json.Unmarshal(body, repo)
	err = repo.Create()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}

func handlePut(w http.ResponseWriter, r *http.Request, repo db.CrudRepository) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	err = repo.Fetch(id)
	if err != nil {
		return
	}
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	json.Unmarshal(body, repo)
	err = repo.Update()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}

func handleDelete(w http.ResponseWriter, r *http.Request, repo db.CrudRepository) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	err = repo.Fetch(id)
	if err != nil {
		return
	}
	err = repo.Delete()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}
