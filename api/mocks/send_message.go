package mocks

import (
	"github.com/stretchr/testify/mock"
)

type SendMessageMock struct {
	mock.Mock
}

func (c *SendMessageMock) SendMessage(magicNumber int) ([]string, error) {
	args := c.Called(magicNumber)

	return args.Get(0).([]string), args.Error(1)
}
