package repository

import (
	"context"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
)

type MessageQueueRepository interface {
	Send(ctx context.Context, message *model.Message) error
	Receive(ctx context.Context, opts *ReceiveMessageOptions) ([]*model.Message, error)
	Delete(ctx context.Context, message *model.Message) error
}

type ReceiveMessageOptions struct {
	MaxNumberOfMessages int32
	WaitTimeSeconds     int32
	VisibilityTimeout   int32
	AttributeNames      []string
}
