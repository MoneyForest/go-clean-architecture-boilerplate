package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSConfig struct {
	Environment string
	Region      string
	Endpoint    string
}

func InitSQS(ctx context.Context, cfg SQSConfig) (*sqs.Client, error) {
	var options []func(*config.LoadOptions) error

	// リージョン設定
	options = append(options, config.WithRegion(cfg.Region))

	// ローカル開発用の設定
	if cfg.Environment == "local" && cfg.Endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: cfg.Endpoint,
			}, nil
		})
		options = append(options, config.WithEndpointResolverWithOptions(customResolver))
	}

	// AWS認証情報とリージョンの設定をロード
	awsCfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// SQSクライアントの作成
	client := sqs.NewFromConfig(awsCfg)

	return client, nil
}
