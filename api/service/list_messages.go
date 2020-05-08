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

// @TODO standardize errors
type ListMessagesError struct {
	reason string
}

func (e ListMessagesError) Error() string {
	return e.reason
}

type listMessages struct {
	repository api.MessageRepository
}

func CreateListMessages(repository api.MessageRepository) ListMessages {
	return listMessages{repository: repository}
}

func (s listMessages) ListMessages(query Query) ([]api.Message, error) {
	var (
		messages []api.Message
		err      error
	)
	messages, err = s.repository.FindByEmail(query.Email)
	if err != nil {
		return messages, err
	}

	return messages, nil
}
