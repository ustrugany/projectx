package service // could be called service

import (
	"fmt"

	"github.com/ustrugany/projectx/api"
)

type SendMessage interface {
	SendMessage(magicNumber int) ([]string, []string, error)
}

type SendMessageError struct {
	Reason string
}

func (e SendMessageError) Error() string {
	return e.Reason
}

type MessageDelivery interface {
	Deliver(message api.Message) error
	BulkDeliver(messages []api.Message) ([]string, []string, error)
}

type sendMessage struct {
	messageRepository api.MessageRepository
	delivery          MessageDelivery
}

func CreateSendMessage(repository api.MessageRepository, delivery MessageDelivery) SendMessage {
	return sendMessage{messageRepository: repository, delivery: delivery}
}

func (s sendMessage) SendMessage(magicNumber int) ([]string, []string, error) {
	var (
		messages []api.Message
		sent     []string
		failed   []string
		err      error
	)

	messages, err = s.messageRepository.FindByMagicNumber(magicNumber)
	if err != nil {
		return sent, failed, SendMessageError{}
	}

	sent, failed, err = s.delivery.BulkDeliver(messages)
	if err != nil {
		return sent, failed, SendMessageError{
			Reason: fmt.Sprintf("delivery failed for %d/%d messages", len(failed), len(messages)),
		}
	}

	count, err := s.messageRepository.DeleteByUUIDs(sent)
	if err != nil {
		return sent, failed, SendMessageError{}
	}

	if count != len(sent) {
		return sent, failed, SendMessageError{
			Reason: fmt.Sprintf("inconsistency error %d sent %d deleted messages", len(sent), count),
		}
	}

	return sent, failed, nil
}
