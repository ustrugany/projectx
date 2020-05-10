package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/ustrugany/projectx/api"
	"github.com/ustrugany/projectx/api/service"
)

type ListMessagesMock struct {
	mock.Mock
}

func (l *ListMessagesMock) ListMessages(query service.ListQuery) ([]api.Message, error) {
	args := l.Called(query)

	return args.Get(0).([]api.Message), args.Error(1)
}
