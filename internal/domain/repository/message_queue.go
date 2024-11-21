package repository

import (
	"context"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
)

type MessageQueueRepository interface {
	SendMessage(ctx context.Context, message *model.Message) error
	ReceiveMessage(ctx context.Context, opts *ReceiveMessageOptions) ([]*model.Message, error)
	DeleteMessage(ctx context.Context, receiptHandle string) error
}

type ReceiveMessageOptions struct {
	MaxNumberOfMessages int32
	WaitTimeSeconds     int32
	VisibilityTimeout   int32
	AttributeNames      []string
}
