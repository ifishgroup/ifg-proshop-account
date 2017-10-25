package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Account struct {
	Db    *sql.DB
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewAccount(db *sql.DB) *Account {
	return &Account{
		Db: db,
	}
}

func (account *Account) Fetch(id int) (err error) {
	err = account.Db.QueryRow(
		"select id, name, email from accounts where id = $1",
		id).Scan(&account.Id, &account.Name, &account.Email)
	return
}

func (account *Account) Create() (err error) {
	statement := "insert into accounts (name, email) values ($1, $2) returning id"
	stmt, err := account.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(account.Name, account.Email).Scan(&account.Id)
	return
}

func (account *Account) Update() (err error) {
	_, err = account.Db.Exec(
		"update accounts set name = $2, email = $3 where id = $1",
		account.Id,
		account.Name,
		account.Email)
	return
}

func (account *Account) Delete() (err error) {
	_, err = account.Db.Exec("delete from accounts where id = $1", account.Id)
	return
}
