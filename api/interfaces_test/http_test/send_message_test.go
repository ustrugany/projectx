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

func TestSendMessage(t *testing.T) {
	body := []byte(`{"magic_number":123}`)
	req, err := http.NewRequest("POST", "/api/send", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Dependencies
	var config api.Config
	config.Schema = map[string]string{
		"send_message": "../../config/schema/http/send_message.json",
	}

	baseLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	logger := *baseLogger.Sugar()
	sendMessage := mocks.SendMessageMock{}

	// Query and mock setup
	magicNumber := 123
	sendMessage.On("SendMessage", magicNumber).
		Return([]string{"uuid-1"}, nil).
		Once()

	rr := httptest.NewRecorder()
	handler := httpInterface.CreateSendMessageHandler(&sendMessage, logger, config)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	sendMessage.AssertExpectations(t)
}
