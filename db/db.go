package db

import (
	"database/sql"
	"fmt"
)

type CrudRepository interface {
	Fetch(id int) (err error)
	Create() (err error)
	Update() (err error)
	Delete() (err error)
}

func NewDbConnection() (*sql.DB, error) {
	return NewDbConnectionWithParams("postgres", "postgres", "postgres")
}

func NewDbConnectionWithParams(host string, user string, password string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s sslmode=disable", host, user, password))
	if err != nil {
		return &sql.DB{}, err
	}
	return db, nil
}
