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

type MessageRepository interface {
	GetMessageByUUID(uuid string) (Message, error)
	FindMessagesByEmail(email string) ([]Message, error)
	CreateMessage(title, content, email string, magicNumber int) (Message, error)
	UpdateMessage(uuid, title, content, email string, magicNumber int) (Message, error)
	DeleteByMagicNumber(magicNumber int) (int, error)
}
