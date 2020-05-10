// +build functional

package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/ustrugany/projectx/api"
	httpInterface "github.com/ustrugany/projectx/api/interfaces/http"
	"github.com/ustrugany/projectx/api/mocks"
	"github.com/ustrugany/projectx/api/service"
)

func TestListMessages(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/messaages", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"email": "test@example.com",
	}
	req = mux.SetURLVars(req, vars)

	// Dependencies
	var config api.Config
	baseLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	logger := *baseLogger.Sugar()
	listMessages := mocks.ListMessagesMock{}

	// Query and mock setup
	query := service.ListQuery{Email: "test@example.com", MagicNumber: "", Title: "", PageSize: 3, PageToken: 0}
	messages := []api.Message{
		{UUID: "uuid-1", Title: "title1", Content: "content1", Email: "test@example.com", MagicNumber: 123},
		{UUID: "uuid-2", Title: "title2", Content: "content2", Email: "test@example.com", MagicNumber: 234},
	}
	listMessages.On("ListMessages", query).Return(messages, nil).Times(1)

	rr := httptest.NewRecorder()
	handler := httpInterface.CreateListMessagesHandler(&listMessages, logger, config)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"data":[{"UUID":"uuid-1","Title":"title1","Content":"content1","Email":"test@example.com","MagicNumber":123},{"UUID":"uuid-2","Title":"title2","Content":"content2","Email":"test@example.com","MagicNumber":234}],"meta":{"next_page_token":0}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestListMessagesWithPagination(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/messaages", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"email": "test@example.com",
	}
	req = mux.SetURLVars(req, vars)
	q := req.URL.Query()
	q.Add("page_size", "3")
	q.Add("page_token", "123")
	req.URL.RawQuery = q.Encode()

	// Dependencies
	var config api.Config
	baseLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	logger := *baseLogger.Sugar()
	listMessages := mocks.ListMessagesMock{}

	// Query and mock setup
	query := service.ListQuery{Email: "test@example.com", MagicNumber: "", Title: "", PageSize: 3, PageToken: 123}
	messages := []api.Message{
		{UUID: "uuid-1", Title: "title1", Content: "content1", Email: "test@example.com", MagicNumber: 123},
		{UUID: "uuid-2", Title: "title2", Content: "content2", Email: "test@example.com", MagicNumber: 234},
		{UUID: "uuid-3", Title: "title3", Content: "content3", Email: "test@example.com", MagicNumber: 345},
	}
	listMessages.On("ListMessages", query).Return(messages, nil).Times(1)

	rr := httptest.NewRecorder()
	handler := httpInterface.CreateListMessagesHandler(&listMessages, logger, config)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"data":[{"UUID":"uuid-1","Title":"title1","Content":"content1","Email":"test@example.com","MagicNumber":123},{"UUID":"uuid-2","Title":"title2","Content":"content2","Email":"test@example.com","MagicNumber":234},{"UUID":"uuid-3","Title":"title3","Content":"content3","Email":"test@example.com","MagicNumber":345}],"meta":{"next_page_token":345}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
