package dto

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs/entity"
)

func FromSQSMessage(sqsMsg types.Message) *entity.Message {
	msg := &entity.Message{
		MessageId:         *sqsMsg.MessageId,
		Body:              *sqsMsg.Body,
		ReceiptHandle:     *sqsMsg.ReceiptHandle,
		Attributes:        make(map[string]string),
		MessageAttributes: make(map[string]entity.MessageAttribute),
	}

	for k, v := range sqsMsg.Attributes {
		msg.Attributes[k] = v
	}

	for k, v := range sqsMsg.MessageAttributes {
		msg.MessageAttributes[k] = entity.MessageAttribute{
			DataType:    *v.DataType,
			StringValue: v.StringValue,
			BinaryValue: v.BinaryValue,
		}
	}

	return msg
}

func ToSQSMessageAttributes(messageAttributes map[string]entity.MessageAttribute) map[string]types.MessageAttributeValue {
	if messageAttributes == nil {
		return nil
	}

	sqsAttrs := make(map[string]types.MessageAttributeValue)
	for k, v := range messageAttributes {
		sqsAttrs[k] = types.MessageAttributeValue{
			DataType:    aws.String(v.DataType),
			StringValue: v.StringValue,
			BinaryValue: v.BinaryValue,
		}
	}
	return sqsAttrs
}
