// +build functional

package http_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/ustrugany/projectx/api"
	httpInterface "github.com/ustrugany/projectx/api/interfaces/http"
	"github.com/ustrugany/projectx/api/mocks"
)

func TestCreateMessage(t *testing.T) {
	body := []byte(`{"title":"title1","content":"content1","email":"test@example.com","magic_number":123}`)
	req, err := http.NewRequest("POST", "/api/messaage", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Dependencies
	var config api.Config
	config.Schema = map[string]string{
		"create_message": "../../config/schema/http/create_message.json",
	}

	baseLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	logger := *baseLogger.Sugar()
	createMessage := mocks.CreateMessageMock{}

	// Query and mock setup
	message := api.Message{UUID: "uuid-1", Title: "title1", Content: "content1", Email: "test@example.com", MagicNumber: 123}
	createMessage.On("CreateMessage", message.Title, message.Content, message.Email, message.MagicNumber).
		Return(message, nil).
		Once()

	rr := httptest.NewRecorder()
	handler := httpInterface.CreateCreateMessageHandler(&createMessage, logger, config)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	createMessage.AssertExpectations(t)
}
