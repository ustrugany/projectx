package service // could be called usecase

import (
	"fmt"

	"github.com/ustrugany/projectx/api"
)

type SendMessage interface {
	// Sends the messages by magicNumber.
	// Returns delivered uuids and/or an error
	SendMessage(magicNumber int) ([]string, error)
}

type SendMessageError struct {
	reason      string
	failed      []string
	deliveryErr error
	deleteErr   error
}

func CreateSendMessageError(reason string, failed []string, deliveryErr, deleteErr error) SendMessageError {
	return SendMessageError{
		reason:      reason,
		failed:      failed,
		deleteErr:   deleteErr,
		deliveryErr: deliveryErr,
	}
}

func (e SendMessageError) Error() string {
	return e.reason
}

type PartialDeliveryError struct {
	reason  string
	failed  []string
	lastErr error
}

func (e PartialDeliveryError) Error() string {
	return e.reason
}

func CreatePartialDeliveryError(reason string, failed []string, lastErr error) PartialDeliveryError {
	return PartialDeliveryError{
		reason:  reason,
		failed:  failed,
		lastErr: lastErr,
	}
}

type MessageDelivery interface {
	Deliver(message api.Message) error
	BulkDeliver(messages []api.Message) ([]string, error)
}

type sendMessage struct {
	repository api.MessageRepository
	delivery   MessageDelivery
}

func CreateSendMessage(repository api.MessageRepository, delivery MessageDelivery) SendMessage {
	return sendMessage{repository: repository, delivery: delivery}
}

func (s sendMessage) SendMessage(magicNumber int) ([]string, error) {
	var (
		delivered []string
	)

	// Lookup
	messages, err := s.repository.FindByMagicNumber(magicNumber)
	if err != nil {
		return delivered, fmt.Errorf("failed to send messages: %w", err)
	}
	if len(messages) == 0 {
		return delivered, nil
	}

	// Delivery, it is not atomic operation, can return partial result
	delivered, err = s.delivery.BulkDeliver(messages)
	if err != nil {
		// Either failed to sent all or unexpected error
		deliveryErr, ok := err.(PartialDeliveryError)
		if !ok || (len(deliveryErr.failed) == len(messages)) {
			return delivered, fmt.Errorf("failed to send messages: %w", err)
		}

		// Bulk delete the ones delivered
		_, deleteErr := s.repository.DeleteByUUIDs(delivered)
		if deleteErr != nil {
			return delivered, fmt.Errorf(
				"failed to complete sending messages: %w", CreateSendMessageError(
					"partially delivered some messages but failed to delete them",
					deliveryErr.failed,
					deliveryErr,
					deleteErr,
				),
			)
		}

		return delivered, nil
	}

	// Bulk delete
	if len(delivered) > 0 {
		_, deleteErr := s.repository.DeleteByUUIDs(delivered)
		if deleteErr != nil {
			return delivered, fmt.Errorf(
				"failed to complete sending messages: %w", CreateSendMessageError(
					"failed to delete delivered messages",
					nil,
					nil,
					deleteErr,
				),
			)
		}
	}

	return delivered, nil
}
