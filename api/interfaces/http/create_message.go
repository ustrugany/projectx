package http

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/qri-io/jsonschema"
	"go.uber.org/zap"

	"github.com/ustrugany/projectx/api"
	"github.com/ustrugany/projectx/api/service"
)

type CreateMessageRequest struct {
	Title       string `json:"title" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Content     string `json:"content" validate:"required,gte=0,lte=1500"`
	MagicNumber int    `json:"magic_number" validate:"required,gte=0,lte=140"`
}

type createMessageHandler struct {
	errorHandler
	service service.CreateMessage
	config  api.Config
}

func CreateCreateMessageHandler(
	service service.CreateMessage,
	logger zap.SugaredLogger,
	config api.Config,
) http.Handler {
	return createMessageHandler{
		service: service,
		config:  config,
		errorHandler: errorHandler{
			logger: logger,
		},
	}
}

func (h createMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// Validate schema
	schemaValidation := struct {
		Errors []jsonschema.ValError `json:"errors"`
	}{}
	if schemaValidation.Errors, err = rs.ValidateBytes(body); len(schemaValidation.Errors) > 0 {
		h.logger.Warnw("message found invalid", "schemaValidation", schemaValidation)
		var bytes []byte
		bytes, err = json.Marshal(schemaValidation)
		if err != nil {
			panic(err)
		}
		h.serve400Error(err, string(bytes), w)
		return
	}

	// Unmarshal dto
	var cm CreateMessageRequest
	err = json.Unmarshal(body, &cm)
	if err != nil {
		h.serve400Error(err, BadRequestBody, w)
		return
	}

	var (
		validationErr service.ValidationError
		message       api.Message
	)

	message, err = h.service.CreateMessage(cm.Title, cm.Content, cm.Email, cm.MagicNumber)
	if err != nil {
		// Something went wrong
		if !errors.As(err, &validationErr) {
			h.serve500Error(err, ServerErrorBody, w)
			return
		}

		// Validation errors
		validation := struct {
			Errors map[string]string `json:"errors"`
		}{
			Errors: validationErr.Violations(),
		}
		var bytes []byte
		bytes, err = json.Marshal(validation)
		if err != nil {
			panic(err)
		}
		h.serve400Error(err, string(bytes), w)
		return
	}

	h.logger.Debugw("creating message", "message", message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
