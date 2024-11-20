package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Key string

// 実際のキューIDとキュー名のマップにすることで環境差異を吸収
const (
	SQSKeySample Key = "sample"
)

type SQSConfig struct {
	Environment string
	Region      string
	Endpoint    string
	QueueNames  map[Key]string
}

type SQSClient struct {
	Client    *sqs.Client
	QueueURLs map[Key]string
}

func InitSQS(ctx context.Context, cfg SQSConfig) (*SQSClient, error) {
	var options []func(*config.LoadOptions) error
	options = append(options, config.WithRegion(cfg.Region))

	switch cfg.Environment {
	case "local", "test":
		if cfg.Endpoint == "" {
			return nil, fmt.Errorf("SQS endpoint is required for local/test environment")
		}
		options = append(options, config.WithCredentialsProvider(aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     "dummy",
					SecretAccessKey: "dummy",
					SessionToken:    "dummy",
				}, nil
			},
		)))
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == sqs.ServiceID {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               cfg.Endpoint,
					SigningRegion:     cfg.Region,
					HostnameImmutable: true,
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown service %s", service)
		})
		options = append(options, config.WithEndpointResolverWithOptions(customResolver))
	default:
		return nil, fmt.Errorf("invalid environment: %s", cfg.Environment)
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := sqs.NewFromConfig(awsCfg)
	sqsClient := &SQSClient{
		Client:    client,
		QueueURLs: make(map[Key]string),
	}

	for Key, queueName := range cfg.QueueNames {
		queueURL := fmt.Sprintf("%s/000000000000/%s", cfg.Endpoint, queueName)
		if cfg.Environment != "local" && cfg.Environment != "test" {
			// 本番環境の場合は GetQueueUrl API を使用
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
