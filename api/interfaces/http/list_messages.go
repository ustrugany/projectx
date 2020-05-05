package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/ustrugany/projectx/api"
	"github.com/ustrugany/projectx/api/service"
)

type ListMessagesHandler struct {
	useCase service.ListMessagesUseCase
	logger  zap.SugaredLogger
	config  api.Config
}

func (h ListMessagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := w.Write([]byte(fmt.Sprintf("Hola %s!", vars["email"])))
	h.logger.Infow("Listing messages", "m", "@TODO put message here")
	if err != nil {
		panic(err)
	}
}

func CreateListMessagesHandler(useCase service.ListMessagesUseCase, logger zap.SugaredLogger, config api.Config) ListMessagesHandler {
	return ListMessagesHandler{
		useCase: useCase,
		logger:  logger,
		config:  config,
	}
}
