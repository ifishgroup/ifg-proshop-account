package db

type FakeAccount struct {
	Id    int
	Name  string
	Email string
}

func (account *FakeAccount) Fetch(id int) (err error) {
	account.Id = id
	return
}

func (account *FakeAccount) Create() (err error) {
	return
}

func (account *FakeAccount) Update() (err error) {
	return
}

func (account *FakeAccount) Delete() (err error) {
	return
}
