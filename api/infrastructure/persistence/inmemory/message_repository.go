package inmemory

import (
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/ustrugany/projectx/api"
)

type MessageRepository struct {
	db DB
}

func CreateMessageRepository(db DB) MessageRepository {
	return MessageRepository{db: db}
}

func (mr MessageRepository) Create(title, content, email string, magicNumber int) (api.Message, error) {
	u := uuid.NewV4()
	message := api.NewMessage(u.String(), title, content, email, magicNumber)
	if err := mr.db.Save(message); err != nil {
		return message, err
	}

	return message, nil
}

func (mr MessageRepository) GetByUUID(uuid string) (api.Message, error) {
	var (
		message api.Message
		err     error
	)
	if message, err = mr.db.Get(uuid); err != nil {
		return message, err
	}

	return message, nil
}

func (mr MessageRepository) FindByEmail(email string, pageSize, pageToken int) ([]api.Message, error) {
	messages, err := mr.db.GetByEmail(email)
	if err != nil {
		return []api.Message{}, fmt.Errorf("failed to find messages: %w", err)
	}

	return messages, nil
}

func (mr MessageRepository) FindByMagicNumber(magicNumber int) ([]api.Message, error) {
	messages, err := mr.db.GetByMagicNumber(magicNumber)
	if err != nil {
		return []api.Message{}, fmt.Errorf("failed to find messages: %w", err)
	}

	return messages, nil
}

func (mr MessageRepository) Update(uuid, title, content, email string, magicNumber int) (api.Message, error) {
	message, err := mr.GetByUUID(uuid)
	if err != nil {
		return message, err
	}

	message.Content = content
	message.Title = title
	message.Email = email
	message.MagicNumber = magicNumber

	if err := mr.db.Save(message); err != nil {
		return message, err
	}

	return message, nil
}

func (mr MessageRepository) DeleteByUUID(uuid string) (int, error) {
	var (
		err   error
		count int
	)

	if count, err = mr.db.Delete(uuid); err != nil {
		return 0, err
	}

	return count, nil
}

func (mr MessageRepository) DeleteByUUIDs(uuids []string) (int, error) {
	var total int
	for _, value := range uuids {
		if count, err := mr.db.Delete(value); err == nil {
			total += count
		}
	}

	return total, nil
}
