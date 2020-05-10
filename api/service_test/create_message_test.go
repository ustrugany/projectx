// +build unit

package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ustrugany/projectx/api"
	"github.com/ustrugany/projectx/api/mocks"
	"github.com/ustrugany/projectx/api/service"
)

func TestCreateMessageSuccess(t *testing.T) {
	repository := &mocks.MessageRepositoryMock{}
	validator := &mocks.MessageValidatorMock{}
	type testCase struct {
		title       string
		email       string
		content     string
		magicNumber int
		violations  map[string]string
		message     api.Message
	}
	var testCases = []testCase{
		{
			title:       "test1",
			email:       "test1@test1.com",
			content:     "test1",
			magicNumber: 1,
			violations:  map[string]string{},
			message:     api.NewMessage("", "test1", "test1@test1.com", "test1", 1),
		},
		{
			title:       "test1",
			email:       "test1%test1.com",
			content:     "test1",
			magicNumber: 1,
			violations:  map[string]string{},
			message:     api.Message{},
		},
	}
	s := service.CreateCreateMessage(repository, validator)
	for _, tc := range testCases {
		repository.On("Create", tc.title, tc.content, tc.email, tc.magicNumber).
			Return(tc.message, nil)
		validator.On("ValidateForCreate", tc.title, tc.content, tc.email, tc.magicNumber).
			Return(tc.violations, nil)
		actual, err := s.CreateMessage(tc.title, tc.content, tc.email, tc.magicNumber)
		assert.NoError(t, err)
		assert.EqualValues(t, tc.message, actual, "expected %v, actual %v", tc.message, actual)
		repository.AssertExpectations(t)
		validator.AssertExpectations(t)
	}
}

func TestCreateMessageFailedValidation(t *testing.T) {
	repository := &mocks.MessageRepositoryMock{}
	validator := &mocks.MessageValidatorMock{}
	type testCase struct {
		title       string
		email       string
		content     string
		magicNumber int
		violations  map[string]string
		message     api.Message
		err         error
	}
	var testCases = []testCase{
		{
			title:       "test1",
			email:       "test1%test1.com",
			content:     "test1",
			magicNumber: 1,
			violations:  map[string]string{"email": "invalid email"},
			message:     api.Message{},
			err:         service.ValidationError{},
		},
	}
	s := service.CreateCreateMessage(repository, validator)
	for _, tc := range testCases {
		validator.On("ValidateForCreate", tc.title, tc.content, tc.email, tc.magicNumber).
			Return(tc.violations, nil).Once()
		actual, err := s.CreateMessage(tc.title, tc.content, tc.email, tc.magicNumber)
		assert.Error(t, err)
		if assert.IsType(t, tc.err, err) {
			validationErr, _ := err.(service.ValidationError)
			assert.Equal(t, validationErr.Violations(), tc.violations)
		}
		assert.EqualValues(t, tc.message, actual, "expected %v, actual %v", tc.message, actual)
		repository.AssertExpectations(t)
		validator.AssertExpectations(t)
	}
}

func TestCreateMessageRepositoryError(t *testing.T) {
	repository := &mocks.MessageRepositoryMock{}
	validator := &mocks.MessageValidatorMock{}
	type testCase struct {
		title       string
		email       string
		content     string
		magicNumber int
		message     api.Message
	}
	var testCases = []testCase{
		{
			title:       "test1",
			email:       "test1%test1.com",
			content:     "test1",
			magicNumber: 1,
			message:     api.Message{},
		},
	}
	s := service.CreateCreateMessage(repository, validator)
	for _, tc := range testCases {
		repository.On("Create", tc.title, tc.content, tc.email, tc.magicNumber).
			Return(tc.message, errors.New("some error")).Once()
		validator.On("ValidateForCreate", tc.title, tc.content, tc.email, tc.magicNumber).
			Return(map[string]string{}, nil).Once()
		actual, err := s.CreateMessage(tc.title, tc.content, tc.email, tc.magicNumber)
		assert.Error(t, err)
		assert.EqualValues(t, tc.message, actual, "expected %v, actual %v", tc.message, actual)
		repository.AssertExpectations(t)
		validator.AssertExpectations(t)
	}
}

func TestCreateMessageValidatorError(t *testing.T) {
	repository := &mocks.MessageRepositoryMock{}
	validator := &mocks.MessageValidatorMock{}
	type testCase struct {
		title       string
		email       string
		content     string
		magicNumber int
		message     api.Message
	}
	var testCases = []testCase{
		{
			title:       "test1",
			email:       "test1%test1.com",
			content:     "test1",
			magicNumber: 1,
			message:     api.Message{},
		},
	}
	s := service.CreateCreateMessage(repository, validator)
	for _, tc := range testCases {
		validator.On("ValidateForCreate", tc.title, tc.content, tc.email, tc.magicNumber).
			Return(map[string]string{}, errors.New("some error")).Once()
		actual, err := s.CreateMessage(tc.title, tc.content, tc.email, tc.magicNumber)
		assert.Error(t, err)
		assert.EqualValues(t, tc.message, actual, "expected %v, actual %v", tc.message, actual)
		repository.AssertNotCalled(t, "Create")
		validator.AssertExpectations(t)
	}
}
