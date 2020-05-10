package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/ustrugany/projectx/api"
	"github.com/ustrugany/projectx/api/service"
)

type Meta struct {
	NextPageToken int `json:"next_page_token"`
}

type ListMessagesResponse struct {
	Data []api.Message `json:"data"`
	Meta Meta          `json:"meta"`
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
	defaultPageSize, defaultPageToken := 3, 0
	pageSize, pageToken := defaultPageSize, defaultPageToken
	if email, ok := vars["email"]; ok {
		query := r.URL.Query()
		pageSizeParam := query.Get("page_size")
		if pageSizeParam != "" {
			if tmp, err := strconv.Atoi(pageSizeParam); err == nil {
				pageSize = tmp
			}
		}

		pageTokenParam := query.Get("page_token")
		if pageTokenParam != "" {
			if tmp, err := strconv.Atoi(pageTokenParam); err == nil {
				pageToken = tmp
			}
		}

		messages, err := h.service.ListMessages(service.ListQuery{
			Email:     email,
			PageToken: pageToken,
			PageSize:  pageSize,
		})
		if err != nil {
			h.serve500Error(err, ServerErrorBody, w)
		}
		h.logger.Debugw("listing messages", "email", email, "messages_count", len(messages))

		response := ListMessagesResponse{
			Data: messages,
		}
		if len(messages) == pageSize {
			lastSeenMessage := messages[pageSize-1]
			response.Meta = Meta{
				NextPageToken: lastSeenMessage.MagicNumber,
			}
		}

		bytes, err := json.Marshal(response)
		if err != nil {
			h.serve500Error(err, ServerErrorBody, w)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(bytes)
		if err != nil {
			h.serve500Error(err, ServerErrorBody, w)
		}

		return
	}

	h.serve400Error(nil, BadRequestBody, w)
}
