package service // could be called service

import (
	"fmt"
	"strings"

	"github.com/ustrugany/projectx/api"
)

type CreateMessage interface {
	CreateMessage(title, content, email string, magicNumber int)
}

type MessageValidationError struct {
	Reason string
}

func (e MessageValidationError) Error() string {
	return e.Reason
}

type CreateMessageUseCase struct {
	repository api.MessageRepository
}

func CreateCreateMessageUseCase(repository api.MessageRepository) CreateMessageUseCase {
	return CreateMessageUseCase{repository: repository}
}

func (uc CreateMessageUseCase) isEmailValid(email string) bool {
	return strings.Contains(email, "@")
}

func (uc CreateMessageUseCase) CreateMessage(title, content, email string, magicNumber int) (api.Message, error) {
	var (
		message api.Message
		err     error
	)
	if !uc.isEmailValid(email) {
		return message, MessageValidationError{Reason: fmt.Sprintf("invalid email %s provided", email)}
	}

	message, err = uc.repository.CreateMessage(title, content, email, magicNumber)
	if err != nil {
		return api.Message{}, fmt.Errorf("failed to create message: %w", err)
	}

	return message, nil
}
