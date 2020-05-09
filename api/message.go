package api

import (
	"html"
	"strings"
)

type Message struct {
	UUID        string
	Title       string
	Content     string
	Email       string
	MagicNumber int
}

func NewMessage(uuid, title, content, email string, magicNumber int) Message {
	return Message{
		UUID:        uuid,
		Title:       html.EscapeString(strings.TrimSpace(title)),
		Content:     html.EscapeString(strings.TrimSpace(content)),
		Email:       email,
		MagicNumber: magicNumber,
	}
}

func (m Message) Validate() (bool, map[string]string) {
	errors := make(map[string]string)

	if m.UUID == "" {
		errors["UUID"] = "unique identifier required"
	}
	if m.Title == "" {
		errors["Title"] = "title required"
	}
	if m.Content == "" || len(m.Content) > 500 {
		errors["Content"] = "content length should be between 0 and 500 characters"
	}
	if m.Email == "" {
		errors["Email"] = "valid email required"
	}

	return len(errors) == 0, errors
}

type RepositoryError struct {
	reason string
}

func (e RepositoryError) Error() string {
	return e.reason
}

type MessageRepository interface {
	Create(title, content, email string, magicNumber int) (Message, error)
	GetByUUID(uuid string) (Message, error)
	FindByEmail(email string, pageSize, pageToken int) ([]Message, error)
	FindByMagicNumber(magicNumber int) ([]Message, error)
	Update(uuid, title, content, email string, magicNumber int) (Message, error)
	DeleteByUUID(uuid string) (int, error)
	DeleteByUUIDs(uuids []string) (int, error)
}
