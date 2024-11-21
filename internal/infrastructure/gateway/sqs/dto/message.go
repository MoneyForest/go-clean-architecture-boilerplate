package dto

import (
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/sqs/entity"
)

func ToMessageModel(entity *entity.Message) *model.Message {
	if entity == nil {
		return nil
	}

	messageAttributes := make(map[string]model.MessageAttribute, len(entity.MessageAttributes))
	for key, value := range entity.MessageAttributes {
		messageAttributes[key] = model.MessageAttribute{
			DataType:    value.DataType,
			StringValue: value.StringValue,
			BinaryValue: value.BinaryValue,
		}
	}

	return &model.Message{
		MessageId:         entity.MessageId,
		Body:              entity.Body,
		ReceiptHandle:     entity.ReceiptHandle,
		Attributes:        entity.Attributes,
		MessageAttributes: messageAttributes,
	}
}

func ToMessageEntity(model *model.Message) *entity.Message {
	if model == nil {
		return nil
	}

	messageAttributes := make(map[string]entity.MessageAttribute, len(model.MessageAttributes))
	for key, value := range model.MessageAttributes {
		messageAttributes[key] = entity.MessageAttribute{
			DataType:    value.DataType,
			StringValue: value.StringValue,
			BinaryValue: value.BinaryValue,
		}
	}

	return &entity.Message{
		MessageId:         model.MessageId,
		Body:              model.Body,
		ReceiptHandle:     model.ReceiptHandle,
		Attributes:        model.Attributes,
		MessageAttributes: messageAttributes,
	}
}
