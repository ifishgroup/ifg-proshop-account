package main

import (
	"github.com/ifishgroup/ifg-proshop-account/db"
	"github.com/ifishgroup/ifg-proshop-account/server"
)

func main() {
	db, err := db.NewDbConnection()
	if err != nil {
		panic(err)
	}
	server := server.NewServer(db)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
