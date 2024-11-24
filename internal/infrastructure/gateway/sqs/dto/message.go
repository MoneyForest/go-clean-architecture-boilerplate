package dto

import (
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
)

type MessageEntity struct {
	Body              string
	MessageAttributes map[string]string
}

func ToMessageEntity(message *model.Message) *MessageEntity {
	return &MessageEntity{
		Body:              message.Body,
		MessageAttributes: message.Attributes,
	}
}

func ToMessageModel(entity *MessageEntity) *model.Message {
	return model.NewMessage(entity.Body, entity.MessageAttributes)
}
