package service // could be called service

import (
	"fmt"
	"strings"

	"github.com/ustrugany/projectx/api"
)

type UpdateMessage interface {
	UpdateMessage(uuid, title, content, email string, magicNumber int) (api.Message, error)
}

type UpdateMessageError struct {
	reason string
}

func (e UpdateMessageError) Error() string {
	return e.reason
}

type updateMessageUseCase struct {
	repository api.MessageRepository
}

func CreateUpdateMessageUseCase(repository api.MessageRepository) UpdateMessage {
	return updateMessageUseCase{repository: repository}
}

func (uc updateMessageUseCase) isEmailValid(email string) bool {
	return strings.Contains(email, "@")
}

func (uc updateMessageUseCase) UpdateMessage(uuid, title, content, email string, magicNumber int) (api.Message, error) {
	var (
		message api.Message
		err     error
	)

	message, err = uc.repository.GetByUUID(uuid)
	if err != nil {
		return message, fmt.Errorf("failed to find message with uuid %s: %w", uuid, err)
	}
	if message.UUID == "" {
		return message, fmt.Errorf("message with uuid %s: not found", uuid)
	}

	if !uc.isEmailValid(email) {
		return message, UpdateMessageError{reason: fmt.Sprintf("invalid email %s provided", email)}
	}

	message, err = uc.repository.Update(uuid, title, content, email, magicNumber)
	if err != nil {
		return api.Message{}, fmt.Errorf("failed to create message: %w", err)
	}

	return message, nil
}
