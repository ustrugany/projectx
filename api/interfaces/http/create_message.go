package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/qri-io/jsonschema"
	"go.uber.org/zap"

	"github.com/ustrugany/projectx/api"
	"github.com/ustrugany/projectx/api/service"
)

type CreateMessageRequest struct {
	Title       string `json:"title"`
	Email       string `json:"email"`
	Content     string `json:"content"`
	MagicNumber int    `json:"magic_number"`
}

type CreateMessageHandler struct {
	service service.CreateMessage
	logger  zap.SugaredLogger
	config  api.Config
}

func CreateCreateMessageHandler(
	service service.CreateMessage,
	logger zap.SugaredLogger,
	config api.Config,
) http.Handler {
	return CreateMessageHandler{
		service: service,
		logger:  logger,
		config:  config,
	}
}

func (h CreateMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.serve500Error(err, ServerErrorBody, w)
		return
	}

	// Read the schema file
	path, err := h.config.SchemaPath("create_message")
	if err != nil {
		h.serve400Error(err, BadRequestBody, w)
		return
	}
	fileSchema, err := ioutil.ReadFile(path)
	if err != nil {
		h.serve400Error(err, BadRequestBody, w)
		return
	}
	rs := &jsonschema.RootSchema{}
	if err = json.Unmarshal(fileSchema, rs); err != nil {
		h.serve500Error(err, ServerErrorBody, w)
		return
	}

	// Validate
	validation := struct {
		Errors []jsonschema.ValError `json:"errors"`
	}{}
	if validation.Errors, err = rs.ValidateBytes(body); len(validation.Errors) > 0 {
		h.logger.Warnw("message found invalid", "validation", validation)
		var errorsData []byte
		errorsData, err = json.Marshal(validation)
		if err != nil {
			panic(err)
		}
		h.serve400Error(err, string(errorsData), w)
		return
	}

	// Unmarshal dto
	var cm CreateMessageRequest
	err = json.Unmarshal(body, &cm)
	if err != nil {
		h.serve400Error(err, BadRequestBody, w)
		return
	}

	var m api.Message
	m, err = h.service.CreateMessage(cm.Title, cm.Content, cm.Email, cm.MagicNumber)
	if err != nil {
		h.serve500Error(err, ServerErrorBody, w)
		return
	}
	h.logger.Debugw("message created", "message", m)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h CreateMessageHandler) serve500Error(err error, content string, w http.ResponseWriter) {
	h.logger.Error("Server error", "err", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if _, err = fmt.Fprintln(w, content); err != nil {
		panic(err)
	}
}

func (h CreateMessageHandler) serve400Error(err error, content string, w http.ResponseWriter) {
	h.logger.Error("Bad request", "err", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	if _, err = fmt.Fprintln(w, content); err != nil {
		panic(err)
	}
}
