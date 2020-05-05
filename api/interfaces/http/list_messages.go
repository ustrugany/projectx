package http

import (
	"encoding/json"
	"fmt"
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
	service service.ListMessages
	logger  zap.SugaredLogger
	config  api.Config
}

func (h listMessagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if email, ok := vars["email"]; ok {
		messages, err := h.service.ListMessages(service.Query{Email: email})
		if err != nil {
			h.serve500Error(err, ServerErrorBody, w)
		}

		h.logger.Debugw("messages found", "email", email, "messages_count", len(messages))

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

func CreateListMessagesHandler(service service.ListMessages, logger zap.SugaredLogger, config api.Config) http.Handler {
	return listMessagesHandler{
		service: service,
		logger:  logger,
		config:  config,
	}
}

func (h listMessagesHandler) serve500Error(err error, content string, w http.ResponseWriter) {
	h.logger.Error("Server error", "err", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if _, err = fmt.Fprintln(w, content); err != nil {
		panic(err)
	}
}

func (h listMessagesHandler) serve400Error(err error, content string, w http.ResponseWriter) {
	h.logger.Error("Bad request", "err", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	if _, err = fmt.Fprintln(w, content); err != nil {
		panic(err)
	}
}
