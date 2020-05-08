package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/ustrugany/projectx/api"
)

type MessageRepositoryMock struct {
	mock.Mock
}

func (r *MessageRepositoryMock) Create(title, content, email string, magicNumber int) (api.Message, error) {
	args := r.Called(title, content, email, magicNumber)
	return args.Get(0).(api.Message), args.Error(1)
}
func (r *MessageRepositoryMock) GetByUUID(uuid string) (api.Message, error) {
	args := r.Called(uuid)
	return args.Get(0).(api.Message), args.Error(1)
}
func (r *MessageRepositoryMock) FindByEmail(email string) ([]api.Message, error) {
	args := r.Called(email)
	return args.Get(0).([]api.Message), args.Error(1)
}
func (r *MessageRepositoryMock) FindByMagicNumber(magicNumber int) ([]api.Message, error) {
	args := r.Called(magicNumber)
	return args.Get(0).([]api.Message), args.Error(1)
}
func (r *MessageRepositoryMock) Update(uuid, title, content, email string, magicNumber int) (api.Message, error) {
	args := r.Called(uuid, title, content, email, magicNumber)
	return args.Get(0).(api.Message), args.Error(1)
}
func (r *MessageRepositoryMock) DeleteByUUID(uuid string) (int, error) {
	args := r.Called(uuid)
	return args.Get(0).(int), args.Error(1)
}
func (r *MessageRepositoryMock) DeleteByUUIDs(uuids []string) (int, error) {
	args := r.Called(uuids)
	return args.Get(0).(int), args.Error(1)
}
