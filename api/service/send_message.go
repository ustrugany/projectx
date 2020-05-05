package service // could be called service

import (
	"github.com/ustrugany/projectx/api"
)

type SendMessageError struct{}

func (e SendMessageError) Error() string {
	return "Error while sending messages"
}

type SendMessageUseCase struct {
	repository api.MessageRepository
}

func CreateSendMessageUseCase(repository api.MessageRepository) SendMessageUseCase {
	return SendMessageUseCase{repository: repository}
}

func (uc SendMessageUseCase) SendMessage(magicNumber int) (int, error) {
	count, err := uc.repository.DeleteByMagicNumber(magicNumber)
	if err != nil {
		return 0, SendMessageError{}
	}

	return count, SendMessageError{}
}
