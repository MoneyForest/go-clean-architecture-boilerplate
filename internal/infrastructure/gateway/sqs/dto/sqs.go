package dto

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// ToSQSMessageAttributes はメッセージ属性をSQS形式に変換します
func ToSQSMessageAttributes(attributes map[string]string) map[string]types.MessageAttributeValue {
	sqsAttributes := make(map[string]types.MessageAttributeValue)
	for key, value := range attributes {
		sqsAttributes[key] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}
	}
	return sqsAttributes
}

// FromSQSMessage はSQSメッセージからエンティティに変換します
func FromSQSMessage(msg types.Message) *MessageEntity {
	attributes := make(map[string]string)
	for key, attr := range msg.MessageAttributes {
		if attr.StringValue != nil {
			attributes[key] = *attr.StringValue
		}
	}

	return &MessageEntity{
		Body:              *msg.Body,
		MessageAttributes: attributes,
	}
}
