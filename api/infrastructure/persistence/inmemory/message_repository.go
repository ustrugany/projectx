package inmemory

import (
	uuid "github.com/satori/go.uuid"

	"github.com/ustrugany/projectx/api"
)

type MessageRepository struct {
	db DB
}

func CreateMessageRepository(db DB) MessageRepository {
	return MessageRepository{db: db}
}

func (mr MessageRepository) GetMessageByUUID(uuid string) (api.Message, error) {
	var (
		message api.Message
		err     error
	)
	if message, err = mr.db.Get(uuid); err != nil {
		return message, err
	}

	return message, nil
}

func (mr MessageRepository) FindMessagesByEmail(email string) ([]api.Message, error) {
	return []api.Message{}, nil
}

func (mr MessageRepository) CreateMessage(title, content, email string, magicNumber int) (api.Message, error) {
	u := uuid.NewV4()
	message := api.NewMessage(u.String(), title, content, email, magicNumber)
	if err := mr.db.Save(message); err != nil {
		return message, err
	}

	return message, nil
}

func (mr MessageRepository) UpdateMessage(uuid, title, content, email string, magicNumber int) (api.Message, error) {
	message, err := mr.GetMessageByUUID(uuid)
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
