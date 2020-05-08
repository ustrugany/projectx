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
		return errors.New("uuid is required")
	}
	db.Storage[m.UUID] = m

	return nil
}

func (db DB) Get(uuid string) (api.Message, error) {
	var message api.Message
	if len(uuid) == 0 {
		return message, errors.New("uuid is required")
	}

	if message, ok := db.Storage[uuid]; ok {
		return message, nil
	}

	return message, nil
}

func (db DB) GetByEmail(email string) ([]api.Message, error) {
	if len(email) == 0 {
		return []api.Message{}, errors.New("email is required")
	}

	var messages []api.Message
	for uuid, value := range db.Storage {
		if value.Email == email {
			messages = append(messages, db.Storage[uuid])
		}
	}

	if len(messages) == 0 {
		return []api.Message{}, nil
	}

	return messages, nil
}

func (db DB) GetByMagicNumber(magicNumber int) ([]api.Message, error) {
	if magicNumber <= 0 {
		return []api.Message{}, errors.New("magicNumber needs to be positive integer type")
	}

	var messages []api.Message
	for uuid, value := range db.Storage {
		if value.MagicNumber == magicNumber {
			messages = append(messages, db.Storage[uuid])
		}
	}

	if len(messages) == 0 {
		return []api.Message{}, nil
	}

	return messages, nil
}

func (db DB) Delete(uuid string) (int, error) {
	if len(uuid) == 0 {
		return 0, errors.New("uuid is required")
	}

	delete(db.Storage, uuid)

	return 1, nil
}

func CreateDb() DB {
	inMemoryDB := DB{
		Storage: make(map[string]api.Message),
	}

	return inMemoryDB
}
