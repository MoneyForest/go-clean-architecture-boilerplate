package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AWSConfig struct {
	Environment string
	Region      string
	Endpoint    string
}

func LoadAWSConfig(ctx context.Context, cfg AWSConfig) (aws.Config, error) {
	var options []func(*config.LoadOptions) error
	options = append(options, config.WithRegion(cfg.Region))

	switch cfg.Environment {
	case "local", "test":
		if cfg.Endpoint == "" {
			return aws.Config{}, fmt.Errorf("SQS endpoint is required for local/test environment")
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
	default:
		return aws.Config{}, fmt.Errorf("invalid environment: %s", cfg.Environment)
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return awsCfg, nil
}
