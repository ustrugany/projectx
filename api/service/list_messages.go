package service // could be called service

import (
	"github.com/ustrugany/projectx/api"
)

type Query struct {
	Email       string
	MagicNumber string
	Title       string
}

type ListMessagesError struct{}

func (e ListMessagesError) Error() string {
	return "Error while listing messages"
}

type ListMessagesUseCase struct {
	repository api.MessageRepository
}

func CreateListMessagesUseCase(repository api.MessageRepository) ListMessagesUseCase {
	return ListMessagesUseCase{repository: repository}
}

func (uc ListMessagesUseCase) ListMessages(query Query) ([]api.Message, error) {
	var (
		messages []api.Message
		err      error
	)

	messages, err = uc.repository.FindMessagesByEmail(query.Email)
	if err != nil {
		return messages, err
	}

	return messages, nil
}
