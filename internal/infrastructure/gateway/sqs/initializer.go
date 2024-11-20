package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	awscfg "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/aws"
)

type Key string

// 実際のキューIDとキュー名のマップにすることで環境差異を吸収
const (
	SQSKeySample Key = "sample"
)

type SQSConfig struct {
	QueueNames map[Key]string
}

type SQSClient struct {
	Client    *sqs.Client
	QueueURLs map[Key]string
}

func InitSQS(ctx context.Context, awsCfg awscfg.AWSConfig, cfg SQSConfig) (*SQSClient, error) {
	awsConfig, err := awscfg.LoadAWSConfig(ctx, awscfg.AWSConfig{
		Environment: awsCfg.Environment,
		Region:      awsCfg.Region,
		Endpoint:    awsCfg.Endpoint,
	})
	if err != nil {
		return nil, err
	}

	client := sqs.NewFromConfig(awsConfig, func(o *sqs.Options) {
		switch awsCfg.Environment {
		case "local", "test":
			o.BaseEndpoint = aws.String(awsCfg.Endpoint)
		}
	})

	sqsClient := &SQSClient{
		Client:    client,
		QueueURLs: make(map[Key]string),
	}

	for Key, queueName := range cfg.QueueNames {
		var queueURL string
		switch awsCfg.Environment {
		case "local", "test":
			queueURL = fmt.Sprintf("%s/000000000000/%s", awsCfg.Endpoint, queueName)
		default:
			result, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
				QueueName: aws.String(queueName),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to get queue URL for %s: %w", queueName, err)
			}
			queueURL = *result.QueueUrl
		}
		sqsClient.QueueURLs[Key] = queueURL
	}

	return sqsClient, nil
}
