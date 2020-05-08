package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/qri-io/jsonschema"
	"go.uber.org/zap"

	"github.com/ustrugany/projectx/api"
	"github.com/ustrugany/projectx/api/service"
)

const (
	BadRequestBody  = `{"error": "Bad request"}`
	ServerErrorBody = `{"error": "Server error"}`
)

type SendMessageRequest struct {
	MagicNumber int `json:"magic_number"`
}

type sendMessageHandler struct {
	errorHandler
	service service.SendMessage
	config  api.Config
}

func CreateSendMessageHandler(service service.SendMessage, logger zap.SugaredLogger, config api.Config) http.Handler {
	return sendMessageHandler{
		service: service,
		config:  config,
		errorHandler: errorHandler{
			logger: logger,
		},
	}
}

func (h sendMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.serve500Error(err, ServerErrorBody, w)
		return
	}

	// Read the schema file
	path, err := h.config.SchemaPath("send_message")
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
		var bytes []byte
		bytes, err = json.Marshal(validation)
		if err != nil {
			panic(err)
		}
		h.serve400Error(err, string(bytes), w)
		return
	}

	// Unmarshal dto
	var sm SendMessageRequest
	err = json.Unmarshal(body, &sm)
	if err != nil {
		h.serve400Error(err, BadRequestBody, w)
		return
	}

	// Send messages
	sent, err := h.service.SendMessage(sm.MagicNumber)
	if err != nil {
		h.serve500Error(err, ServerErrorBody, w)
		return
	}

	h.logger.Debugw("sending message", "magic_number", sm.MagicNumber, "sent", sent, "failed", failed)

	w.Header().Set("Content-Type", "application/json")
	if len(sent) == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}
