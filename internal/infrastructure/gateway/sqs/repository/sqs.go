package repository

import (
	"context"

	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs/dto"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs/entity"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type ReceiveMessageOptions struct {
	MaxNumberOfMessages int32
	WaitTimeSeconds     int32
	VisibilityTimeout   int32
	AttributeNames      []string
}

type SQSRepository struct {
	sqs       *sqs.Client
	queueName string
}

func NewSQSRepository(sqs *sqs.Client, queueName string) SQSRepository {
	return SQSRepository{sqs: sqs, queueName: queueName}
}

func (r SQSRepository) SendMessage(ctx context.Context, message *entity.Message) error {
	input := &sqs.SendMessageInput{
		QueueUrl:          aws.String(r.queueName),
		MessageBody:       aws.String(message.Body),
		MessageAttributes: dto.ToSQSMessageAttributes(message.MessageAttributes),
	}

	_, err := r.sqs.SendMessage(ctx, input)
	return err
}

func (r SQSRepository) ReceiveMessage(ctx context.Context, opts *ReceiveMessageOptions) ([]*entity.Message, error) {
	if opts == nil {
		opts = &ReceiveMessageOptions{
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
			input.AttributeNames = append(input.AttributeNames, types.QueueAttributeName(attr))
		}
	}

	output, err := r.sqs.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, err
	}

	messages := make([]*entity.Message, len(output.Messages))
	for i, msg := range output.Messages {
		messages[i] = dto.FromSQSMessage(msg)
	}

	return messages, nil
}

func (r SQSRepository) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := r.sqs.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(r.queueName),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}

func (r SQSRepository) ChangeMessageVisibility(ctx context.Context, queueURL string, receiptHandle string, visibilityTimeout int32) error {
	_, err := r.sqs.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(r.queueName),
		ReceiptHandle:     aws.String(receiptHandle),
		VisibilityTimeout: visibilityTimeout,
	})
	return err
}

func (r SQSRepository) SendMessageBatch(ctx context.Context, messages []*entity.Message) error {
	if len(messages) == 0 {
		return nil
	}

	entries := make([]types.SendMessageBatchRequestEntry, len(messages))
	for i, msg := range messages {
		entries[i] = types.SendMessageBatchRequestEntry{
			Id:                aws.String(msg.MessageId),
			MessageBody:       aws.String(msg.Body),
			MessageAttributes: dto.ToSQSMessageAttributes(msg.MessageAttributes),
		}
	}

	input := &sqs.SendMessageBatchInput{
		QueueUrl: aws.String(r.queueName),
		Entries:  entries,
	}

	_, err := r.sqs.SendMessageBatch(ctx, input)
	return err
}

func (r SQSRepository) DeleteMessageBatch(ctx context.Context, messages []*entity.Message) error {
	if len(messages) == 0 {
		return nil
	}

	entries := make([]types.DeleteMessageBatchRequestEntry, len(messages))
	for i, msg := range messages {
		entries[i] = types.DeleteMessageBatchRequestEntry{
			Id:            aws.String(msg.MessageId),
			ReceiptHandle: aws.String(msg.ReceiptHandle),
		}
	}

	input := &sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(r.queueName),
		Entries:  entries,
	}

	_, err := r.sqs.DeleteMessageBatch(ctx, input)
	return err
}
