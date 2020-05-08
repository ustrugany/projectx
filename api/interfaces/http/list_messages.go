package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/ustrugany/projectx/api"
	"github.com/ustrugany/projectx/api/service"
)

type ListMessagesResponse struct {
	Data []api.Message `json:"data"`
}

type listMessagesHandler struct {
	errorHandler
	service service.ListMessages
	config  api.Config
}

func CreateListMessagesHandler(service service.ListMessages, logger zap.SugaredLogger, config api.Config) http.Handler {
	return listMessagesHandler{
		service: service,
		config:  config,
		errorHandler: errorHandler{
			logger: logger,
		},
	}
}

func (h listMessagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if email, ok := vars["email"]; ok {
		messages, err := h.service.ListMessages(service.Query{Email: email})
		if err != nil {
			h.serve500Error(err, ServerErrorBody, w)
		}

		h.logger.Debugw("listing messages", "email", email, "messages_count", len(messages))

		response := ListMessagesResponse{
			Data: messages,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			h.serve500Error(err, ServerErrorBody, w)
		}

		return
	}

	h.serve400Error(nil, BadRequestBody, w)
}
