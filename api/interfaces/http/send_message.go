package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/ustrugany/projectx/api"
	"github.com/ustrugany/projectx/api/service"
)

type SendMessageRequest struct {
	MagicNumber int `json:"magic_number"`
}

type SendMessageHandler struct {
	useCase service.SendMessageUseCase
	logger  zap.SugaredLogger
	config  api.Config
}

func CreateSendMessageHandler(useCase service.SendMessageUseCase, logger zap.SugaredLogger, config api.Config) SendMessageHandler {
	return SendMessageHandler{
		useCase: useCase,
		logger:  logger,
		config:  config,
	}
}
func (h SendMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Sent the message!"))
	h.logger.Infow("Sent a message", "m", "@TODO put message here")
	if err != nil {
		panic(err)
	}
}
