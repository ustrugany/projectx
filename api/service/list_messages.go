package service // could be called service

import (
	"github.com/ustrugany/projectx/api"
)

type ListQuery struct {
	Email       string
	MagicNumber string
	Title       string
	PageSize    int
	PageToken   int
}

type ListMessages interface {
	ListMessages(query ListQuery) ([]api.Message, error)
}

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

func (s listMessages) ListMessages(query ListQuery) ([]api.Message, error) {
	var (
		messages []api.Message
		err      error
	)
	messages, err = s.repository.FindByEmail(query.Email, query.PageSize, query.PageToken)
	if err != nil {
		return messages, err
	}

	return messages, nil
}
