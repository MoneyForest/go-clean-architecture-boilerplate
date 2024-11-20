package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs/entity"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type UserRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Get(ctx context.Context, id uuid.UUID) (*model.User, error)
	List(ctx context.Context, limit, offset int) ([]*model.User, error)
	UpdateTx(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error)
	DeleteTx(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*uuid.UUID, error)
}

type UserCacheRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*model.User, error)
	SetWithTTL(ctx context.Context, user *model.User, ttl time.Duration) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ReceiveMessageOptions struct {
	MaxNumberOfMessages int32
	WaitTimeSeconds     int32
	VisibilityTimeout   int32
	AttributeNames      []string
}

type UserMessageQueueRepository interface {
	SendMessage(ctx context.Context, message *entity.Message) error
	ReceiveMessage(ctx context.Context, opts *ReceiveMessageOptions) ([]*entity.Message, error)
	DeleteMessage(ctx context.Context, receiptHandle string) error
	ChangeMessageVisibility(ctx context.Context, queueURL string, receiptHandle string, visibilityTimeout int32) error
	SendMessageBatch(ctx context.Context, messages []*entity.Message) error
	DeleteMessageBatch(ctx context.Context, messages []*entity.Message) error
}
