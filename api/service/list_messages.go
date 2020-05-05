package service // could be called service

import (
	"github.com/ustrugany/projectx/api"
)

type Query struct {
	Email       string
	MagicNumber string
	Title       string
}

type ListMessages interface {
	ListMessages(query Query) ([]api.Message, error)
}

type ListMessagesError struct {
	Reason string
}

func (e ListMessagesError) Error() string {
	return e.Reason
}

type listMessages struct {
	messageRepository api.MessageRepository
}

func CreateListMessages(repository api.MessageRepository) ListMessages {
	return listMessages{messageRepository: repository}
}

func (s listMessages) ListMessages(query Query) ([]api.Message, error) {
	var (
		messages []api.Message
		err      error
	)
	messages, err = s.messageRepository.FindByEmail(query.Email)
	if err != nil {
		return messages, err
	}

	return messages, nil
}
