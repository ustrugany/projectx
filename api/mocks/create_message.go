package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/ustrugany/projectx/api"
)

type CreateMessageMock struct {
	mock.Mock
}

func (c *CreateMessageMock) CreateMessage(title, content, email string, magicNumber int) (api.Message, error) {
	args := c.Called(title, content, email, magicNumber)

	return args.Get(0).(api.Message), args.Error(1)
}
