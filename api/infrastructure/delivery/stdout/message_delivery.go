package stdout

import (
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

func (ed messageDelivery) BulkDeliver(messages []api.Message) ([]string, []string, error) {
	var (
		err       error
		delivered []string
		failed    []string
	)
	for _, message := range messages {
		ed.logger.Infow("message delivered", "message", message)
		delivered = append(delivered, message.UUID)
	}

	return delivered, failed, err
}
