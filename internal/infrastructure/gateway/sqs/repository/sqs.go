package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/repository"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs/dto"
)

type SQSRepository struct {
	sqs       *sqs.Client
	queueName string
}

func NewSQSRepository(sqs *sqs.Client, queueName string) repository.MessageQueueRepository {
	return &SQSRepository{sqs: sqs, queueName: queueName}
}

func (r *SQSRepository) SendMessage(ctx context.Context, message *model.Message) error {
	messageEntity := dto.ToMessageEntity(message)
	input := &sqs.SendMessageInput{
		QueueUrl:          aws.String(r.queueName),
		MessageBody:       aws.String(messageEntity.Body),
		MessageAttributes: dto.ToSQSMessageAttributes(messageEntity.MessageAttributes),
	}

	_, err := r.sqs.SendMessage(ctx, input)
	return err
}

func (r *SQSRepository) ReceiveMessage(ctx context.Context, opts *repository.ReceiveMessageOptions) ([]*model.Message, error) {
	if opts == nil {
		opts = &repository.ReceiveMessageOptions{
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     20,
		}
	}

	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(r.queueName),
		MaxNumberOfMessages: opts.MaxNumberOfMessages,
		WaitTimeSeconds:     opts.WaitTimeSeconds,
		VisibilityTimeout:   opts.VisibilityTimeout,
	}

	if len(opts.AttributeNames) > 0 {
		for _, attr := range opts.AttributeNames {
			input.MessageSystemAttributeNames = append(input.MessageSystemAttributeNames, types.MessageSystemAttributeName(attr))
		}
	}

	output, err := r.sqs.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, err
	}

	messages := make([]*model.Message, len(output.Messages))
	for i, msg := range output.Messages {
		messageEntity := dto.FromSQSMessage(msg)
		messages[i] = dto.ToMessageModel(messageEntity)
	}

	return messages, nil
}

func (r *SQSRepository) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := r.sqs.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(r.queueName),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}
