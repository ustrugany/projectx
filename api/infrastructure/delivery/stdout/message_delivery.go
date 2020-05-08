package stdout

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/ustrugany/projectx/api"
	"github.com/ustrugany/projectx/api/service"
)

type messageDelivery struct {
	logger zap.SugaredLogger
}

func CreateMessageDelivery(logger zap.SugaredLogger) service.MessageDelivery {
	return messageDelivery{logger: logger}
}

func (ed messageDelivery) Deliver(message api.Message) error {
	ed.logger.Infow("message delivered", "message", message)

	return nil
}

func (ed messageDelivery) BulkDeliver(messages []api.Message) ([]string, error) {
	var (
		err    error
		failed []string
		sent   []string
	)

	if len(messages) == 0 {
		return sent, err
	}

	for _, message := range messages {
		if err := ed.Deliver(message); err != nil {
			failed = append(failed, message.UUID)
			continue
		}
		sent = append(sent, message.UUID)
	}

	if len(failed) > 0 {
		deliveryErr := service.CreatePartialDeliveryError(
			fmt.Sprintf("failed to deliver %d messages", len(failed)),
			failed,
			err,
		)

		return sent, errors.Errorf("delivery error: %w", deliveryErr)
	}

	return sent, err
}
