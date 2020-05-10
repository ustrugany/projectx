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

func TestListMessagesSuccess(t *testing.T) {
	repository := &mocks.MessageRepositoryMock{}
	type testCase struct {
		query   service.ListQuery
		message []api.Message
	}
	var testCases = []testCase{
		{
			query: service.ListQuery{
				Email: "test1@test1.com",
			},
			message: []api.Message{
				api.NewMessage("", "test1", "test1@test1.com", "test1", 1),
				api.NewMessage("", "test1", "test1@test1.com", "test1", 1),
			},
		},
		{
			query: service.ListQuery{
				Email: "test2@test2.com",
			},
			message: []api.Message{
				api.NewMessage("", "test2", "test2@test1.com", "test2", 1),
			},
		},
	}
	s := service.CreateListMessages(repository)
	for _, tc := range testCases {
		repository.On("FindByEmail", tc.query.Email, tc.query.PageSize, tc.query.PageToken).
			Return(tc.message, nil)
		actual, err := s.ListMessages(tc.query)
		assert.NoError(t, err)
		assert.EqualValues(t, tc.message, actual)
		repository.AssertExpectations(t)
	}
}

func TestListMessageRepositoryError(t *testing.T) {
	repository := &mocks.MessageRepositoryMock{}
	type testCase struct {
		query   service.ListQuery
		message []api.Message
	}
	var testCases = []testCase{
		{
			query: service.ListQuery{
				Email: "test1@test1.com",
			},
			message: []api.Message{},
		},
	}
	s := service.CreateListMessages(repository)
	for _, tc := range testCases {
		repository.On("FindByEmail", tc.query.Email, tc.query.PageSize, tc.query.PageToken).
			Return(tc.message, errors.New("some error"))
		_, err := s.ListMessages(tc.query)
		assert.Error(t, err)
		repository.AssertExpectations(t)
	}
}
