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

func TestSendMessageSuccess(t *testing.T) {
	type testCase struct {
		scenario    string
		magicNumber int
		delivered   []string
		messages    []api.Message
	}
	var testCases = []testCase{
		{
			scenario:    "no data",
			magicNumber: 123,
			messages:    []api.Message{},
		},
		{
			scenario:    "all delivered",
			magicNumber: 123,
			delivered:   []string{"uuid-1", "uuid-2"},
			messages: []api.Message{
				api.NewMessage("uuid-1", "test1", "test1@test1.com", "test1", 123),
				api.NewMessage("uuid-2", "test2", "test2@test2.com", "test2", 123),
			},
		},
		{
			scenario:    "all failed to deliver",
			magicNumber: 234,
			messages: []api.Message{
				api.NewMessage("uuid-3", "test3", "test3@test3.com", "test3", 234),
				api.NewMessage("uuid-4", "test4", "test2@test4.com", "test4", 234),
			},
		},
	}
	for _, tc := range testCases {
		repository := &mocks.MessageRepositoryMock{}
		repository.On("FindByMagicNumber", tc.magicNumber).
			Return(tc.messages, nil).
			Once()

		delivery := &mocks.MessageDeliveryMock{}
		if len(tc.messages) > 0 {
			delivery.On("BulkDeliver", tc.messages).
				Return(tc.delivered, nil).
				Once()
		}

		if len(tc.delivered) > 0 {
			repository.On("DeleteByUUIDs", tc.delivered).
				Return(len(tc.delivered), nil).
				Once()
		}

		s := service.CreateSendMessage(repository, delivery)
		delivered, err := s.SendMessage(tc.magicNumber)

		assert.NoError(t, err)
		assert.EqualValues(t, tc.delivered, delivered)

		repository.AssertExpectations(t)
		delivery.AssertExpectations(t)
	}
}

func TestSendMessagePartialDelivery(t *testing.T) {
	type testCase struct {
		scenario    string
		magicNumber int
		delivered   []string
		messages    []api.Message
	}
	var testCases = []testCase{
		{
			scenario:    "half delivered, half failed",
			magicNumber: 234,
			delivered:   []string{"uuid-3"},
			messages: []api.Message{
				api.NewMessage("uuid-3", "test3", "test3@test3.com", "test3", 234),
				api.NewMessage("uuid-4", "test4", "test2@test4.com", "test4", 234),
			},
		},
	}

	for _, tc := range testCases {
		repository := &mocks.MessageRepositoryMock{}
		repository.On("FindByMagicNumber", tc.magicNumber).
			Return(tc.messages, nil)

		delivery := &mocks.MessageDeliveryMock{}
		delivery.On("BulkDeliver", tc.messages).Return(
			tc.delivered, service.CreatePartialDeliveryError("something went wrong",
				[]string{"uuid-4"},
				errors.New("network partition")),
		)

		repository.On("DeleteByUUIDs", tc.delivered).
			Return(len(tc.delivered), nil)

		s := service.CreateSendMessage(repository, delivery)
		delivered, err := s.SendMessage(tc.magicNumber)

		assert.NoError(t, err)
		assert.EqualValues(t, tc.delivered, delivered)

		repository.AssertNumberOfCalls(t, "DeleteByUUIDs", len(tc.delivered))
		repository.AssertNumberOfCalls(t, "FindByMagicNumber", 1)
		delivery.AssertNumberOfCalls(t, "BulkDeliver", 1)
	}
}

func TestSendMessageDeliveryError(t *testing.T) {
	type testCase struct {
		scenario    string
		magicNumber int
		delivered   []string
		messages    []api.Message
	}
	var testCases = []testCase{
		{
			scenario:    "delivery failed",
			magicNumber: 123,
			messages: []api.Message{
				api.NewMessage("uuid-1", "test1", "test1@test1.com", "test1", 123),
				api.NewMessage("uuid-2", "test2", "test2@test2.com", "test2", 123),
			},
		},
	}
	for _, tc := range testCases {
		repository := &mocks.MessageRepositoryMock{}
		repository.On("FindByMagicNumber", tc.magicNumber).
			Return(tc.messages, nil)

		delivery := &mocks.MessageDeliveryMock{}
		delivery.On("BulkDeliver", tc.messages).
			Return(tc.delivered, errors.New("something went wrong"))
		s := service.CreateSendMessage(repository, delivery)

		delivered, err := s.SendMessage(tc.magicNumber)
		assert.Error(t, err)
		assert.EqualValues(t, tc.delivered, delivered)

		repository.AssertNumberOfCalls(t, "FindByMagicNumber", 1)
		delivery.AssertNumberOfCalls(t, "BulkDeliver", 1)
		repository.AssertNotCalled(t, "DeleteByUUIDs")
	}
}

func TestSendMessageRepositoryFindError(t *testing.T) {
	type testCase struct {
		scenario    string
		magicNumber int
		delivered   []string
		messages    []api.Message
	}
	var testCases = []testCase{
		{
			scenario:    "repository find method error",
			magicNumber: 123,
		},
	}
	for _, tc := range testCases {
		repository := &mocks.MessageRepositoryMock{}
		repository.On("FindByMagicNumber", tc.magicNumber).
			Return(tc.messages, errors.New("something went wrong"))

		delivery := &mocks.MessageDeliveryMock{}
		s := service.CreateSendMessage(repository, delivery)

		delivered, err := s.SendMessage(tc.magicNumber)
		assert.Error(t, err)
		assert.EqualValues(t, tc.delivered, delivered)

		repository.AssertNumberOfCalls(t, "FindByMagicNumber", 1)
		delivery.AssertNotCalled(t, "BulkDeliver")
	}
}

func TestSendMessageRepositoryDeleteError(t *testing.T) {
	type testCase struct {
		scenario    string
		magicNumber int
		delivered   []string
		messages    []api.Message
	}
	var testCases = []testCase{
		{
			scenario:    "repository delete method error",
			magicNumber: 123,
			delivered:   []string{"uuid-1"},
			messages: []api.Message{
				api.NewMessage("uuid-1", "test1", "test1@test1.com", "test1", 123),
				api.NewMessage("uuid-2", "test2", "test2@test2.com", "test2", 123),
			},
		},
	}

	for _, tc := range testCases {
		repository := &mocks.MessageRepositoryMock{}
		repository.On("FindByMagicNumber", tc.magicNumber).
			Return(tc.messages, nil)
		repository.On("DeleteByUUIDs", tc.delivered).
			Return(0, errors.New("something went wrong"))

		delivery := &mocks.MessageDeliveryMock{}
		delivery.On("BulkDeliver", tc.messages).
			Return(tc.delivered, nil)

		s := service.CreateSendMessage(repository, delivery)

		delivered, err := s.SendMessage(tc.magicNumber)
		assert.Error(t, err)
		assert.EqualValues(t, tc.delivered, delivered)

		repository.AssertNumberOfCalls(t, "FindByMagicNumber", 1)
	}
}
