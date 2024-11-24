package model

import (
	"time"

	"github.com/google/uuid"
)

type MessageAttributes map[string]string

type Message struct {
	ID            uuid.UUID
	Body          string
	ReceiptHandle string
	Attributes    MessageAttributes
	CreatedAt     time.Time
}

type MessageDTO struct {
	Body              string
	MessageAttributes MessageAttributes
}

func NewMessage(body string, attributes MessageAttributes) *Message {
	return &Message{
		ID:         uuid.New(),
		Body:       body,
		Attributes: attributes,
		CreatedAt:  time.Now(),
	}
}

func (m *Message) Send() *MessageDTO {
	return &MessageDTO{
		Body:              m.Body,
		MessageAttributes: m.Attributes,
	}
}

func (m *Message) SetReceiptHandle(handle string) {
	m.ReceiptHandle = handle
}

func (m *Message) Delete() string {
	return m.ReceiptHandle
}
