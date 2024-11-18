package dto

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Message struct {
	MessageId         string
	Body              string
	ReceiptHandle     string
	Attributes        map[string]string
	MessageAttributes map[string]MessageAttribute
}

type MessageAttribute struct {
	DataType    string
	StringValue *string
	BinaryValue []byte
}

// FromSQSMessage converts SQS message to DTO
func FromSQSMessage(sqsMsg types.Message) *Message {
	msg := &Message{
		MessageId:         *sqsMsg.MessageId,
		Body:              *sqsMsg.Body,
		ReceiptHandle:     *sqsMsg.ReceiptHandle,
		Attributes:        make(map[string]string),
		MessageAttributes: make(map[string]MessageAttribute),
	}

	for k, v := range sqsMsg.Attributes {
		msg.Attributes[k] = v
	}

	for k, v := range sqsMsg.MessageAttributes {
		msg.MessageAttributes[k] = MessageAttribute{
			DataType:    *v.DataType,
			StringValue: v.StringValue,
			BinaryValue: v.BinaryValue,
		}
	}

	return msg
}

// ToSQSMessageAttributes converts DTO message attributes to SQS message attributes
func (m *Message) ToSQSMessageAttributes() map[string]types.MessageAttributeValue {
	if m.MessageAttributes == nil {
		return nil
	}

	sqsAttrs := make(map[string]types.MessageAttributeValue)
	for k, v := range m.MessageAttributes {
		sqsAttrs[k] = types.MessageAttributeValue{
			DataType:    aws.String(v.DataType),
			StringValue: v.StringValue,
			BinaryValue: v.BinaryValue,
		}
	}
	return sqsAttrs
}
