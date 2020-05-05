package service // could be called service

import (
	"fmt"
	"strings"

	"github.com/ustrugany/projectx/api"
)

type CreateMessage interface {
	CreateMessage(title, content, email string, magicNumber int) (api.Message, error)
}

type MessageValidationError struct {
	Reason string
}

func (e MessageValidationError) Error() string {
	return e.Reason
}

type createMessage struct {
	messageRepository api.MessageRepository
}

func CreateCreateMessage(repository api.MessageRepository) CreateMessage {
	return createMessage{messageRepository: repository}
}

func (s createMessage) isEmailValid(email string) bool {
	return strings.Contains(email, "@")
}

func (s createMessage) CreateMessage(title, content, email string, magicNumber int) (api.Message, error) {
	var (
		message api.Message
		err     error
	)
	if !s.isEmailValid(email) {
		return message, MessageValidationError{Reason: fmt.Sprintf("invalid email %s provided", email)}
	}

	message, err = s.messageRepository.Create(title, content, email, magicNumber)
	if err != nil {
		return api.Message{}, fmt.Errorf("failed to create message: %w", err)
	}

	return message, nil
}
