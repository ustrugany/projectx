package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/ustrugany/projectx/api"
)

type MessageDeliveryMock struct {
	mock.Mock
}

func (v *MessageDeliveryMock) Deliver(message api.Message) error {
	args := v.Called(message)
	return args.Error(1)
}

func (v *MessageDeliveryMock) BulkDeliver(messages []api.Message) ([]string, error) {
	args := v.Called(messages)
	return args.Get(0).([]string), args.Error(1)
}
