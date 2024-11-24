package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/repository"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/sqs/dto"
)

// SQSRepository はAWS SQSを使用したメッセージキューリポジトリの実装です
type SQSRepository struct {
	sqs       *sqs.Client
	queueName string
}

// NewSQSRepository は新しいSQSRepositoryインスタンスを作成します
func NewSQSRepository(sqs *sqs.Client, queueName string) repository.MessageQueueRepository {
	return &SQSRepository{
		sqs:       sqs,
		queueName: queueName,
	}
}

// Send はメッセージをSQSキューに送信します
func (r *SQSRepository) Send(ctx context.Context, message *model.Message) error {
	messageDTO := message.Send()
	input := &sqs.SendMessageInput{
		QueueUrl:          aws.String(r.queueName),
		MessageBody:       aws.String(messageDTO.Body),
		MessageAttributes: dto.ToSQSMessageAttributes(messageDTO.MessageAttributes),
	}

	_, err := r.sqs.SendMessage(ctx, input)
	return err
}

// Receive はSQSキューからメッセージを受信します
func (r *SQSRepository) Receive(ctx context.Context, opts *repository.ReceiveMessageOptions) ([]*model.Message, error) {
	if opts == nil {
		opts = &repository.ReceiveMessageOptions{
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     20,
		}
	}

	input := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(r.queueName),
		MaxNumberOfMessages:   opts.MaxNumberOfMessages,
		WaitTimeSeconds:       opts.WaitTimeSeconds,
		VisibilityTimeout:     opts.VisibilityTimeout,
		MessageAttributeNames: []string{"All"},
	}

	if len(opts.AttributeNames) > 0 {
		for _, attr := range opts.AttributeNames {
			input.MessageSystemAttributeNames = append(
				input.MessageSystemAttributeNames,
				types.MessageSystemAttributeName(attr),
			)
		}
	}

	output, err := r.sqs.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, err
	}

	messages := make([]*model.Message, len(output.Messages))
	for i, msg := range output.Messages {
		message := dto.ToMessageModel(dto.FromSQSMessage(msg))
		message.SetReceiptHandle(*msg.ReceiptHandle)
		messages[i] = message
	}

	return messages, nil
}

// DeleteMessage はSQSキューからメッセージを削除します
func (r *SQSRepository) Delete(ctx context.Context, message *model.Message) error {
	_, err := r.sqs.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(r.queueName),
		ReceiptHandle: aws.String(message.Delete()),
	})
	return err
}
