package inmemory

import (
	"errors"

	"github.com/ustrugany/projectx/api"
)

type DB struct {
	Storage map[string]api.Message
}

func (db DB) Save(m api.Message) error {
	if len(m.UUID) == 0 {
		return errors.New("UUID is required")
	}
	db.Storage[m.UUID] = m

	return nil
}

func (db DB) Get(uuid string) (api.Message, error) {
	var m api.Message
	if len(uuid) == 0 {
		return m, errors.New("UUID is required")
	}
	return db.Storage[m.UUID], nil
}

func CreateDb() DB {
	inMemoryDB := DB{
		Storage: make(map[string]api.Message),
	}

	return inMemoryDB
}
