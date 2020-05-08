package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MessageValidatorMock struct {
	mock.Mock
}

func (v *MessageValidatorMock) ValidateForCreate(title, content, email string, magicNumber int) (map[string]string, error) {
	args := v.Called(title, content, email, magicNumber)
	return args.Get(0).(map[string]string), args.Error(1)
}
