package service // could be called use-case, application

import (
	"fmt"

	"github.com/ustrugany/projectx/api"
)

type MassageValidator interface {
	ValidateForCreate(title, content, email string, magicNumber int) (map[string]string, error)
}

type CreateMessage interface {
	CreateMessage(title, content, email string, magicNumber int) (api.Message, error)
}

type ValidationError struct {
	violations map[string]string
	reason     string
}

func (e ValidationError) Error() string {
	return e.reason
}

func (e ValidationError) Violations() map[string]string {
	return e.violations
}

type createMessage struct {
	repository api.MessageRepository
	validator  MassageValidator
}

func CreateCreateMessage(repository api.MessageRepository, validator MassageValidator) CreateMessage {
	return createMessage{repository: repository, validator: validator}
}

func (s createMessage) CreateMessage(title, content, email string, magicNumber int) (api.Message, error) {
	var (
		message    api.Message
		violations map[string]string
		err        error
	)

	violations, err = s.validator.ValidateForCreate(title, content, email, magicNumber)
	if err != nil {
		return api.Message{}, fmt.Errorf("failed to validate message: %w", err)
	}
	if len(violations) > 0 {
		return api.Message{}, ValidationError{
			reason:     "message is invalid",
			violations: violations,
		}
	}

	message, err = s.repository.Create(title, content, email, magicNumber)
	if err != nil {
		return api.Message{}, fmt.Errorf("failed to create message: %w", err)
	}

	return message, nil
}
