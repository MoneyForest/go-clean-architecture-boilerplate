package testhelper

import (
	"context"
	"database/sql"

	mysqlgateway "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql"
	redisgateway "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/redis"
	sqsgateway "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/redis/go-redis/v9"
)

type Gateway struct {
	MySQLClient *sql.DB
	RedisClient *redis.Client
	SQSClient   *sqs.Client
}

type TestEnvironment struct {
	Environment   string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBDatabase    string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	AWSRegion     string
	SQSEndpoint   string
}

func Setup(ctx context.Context) (*Gateway, error) {
	e := &TestEnvironment{
		Environment:   "test",
		DBHost:        "localhost",
		DBPort:        "3306",
		DBUser:        "root",
		DBPassword:    "password",
		DBDatabase:    "maindb",
		RedisHost:     "localhost",
		RedisPort:     "6379",
		RedisPassword: "password",
		AWSRegion:     "ap-northeast-1",
		SQSEndpoint:   "",
	}

	mysqlClient, err := mysqlgateway.InitDB(ctx, mysqlgateway.DBConfig{
		Environment: e.Environment,
		Host:        e.DBHost,
		Port:        e.DBPort,
		User:        e.DBUser,
		Password:    e.DBPassword,
		DBName:      e.DBDatabase,
	})
	if err != nil {
		return nil, err
	}
	redisClient, err := redisgateway.InitRedis(ctx, redisgateway.RedisConfig{
		Environment: e.Environment,
		Host:        e.RedisHost,
		Port:        e.RedisPort,
		Password:    e.RedisPassword,
	})
	if err != nil {
		return nil, err
	}
	sqsClient, err := sqsgateway.InitSQS(ctx, sqsgateway.SQSConfig{
		Environment: e.Environment,
		Region:      e.AWSRegion,
		Endpoint:    e.SQSEndpoint,
	})
	if err != nil {
		return nil, err
	}

	return &Gateway{
		MySQLClient: mysqlClient,
		RedisClient: redisClient,
		SQSClient:   sqsClient,
	}, nil
}
