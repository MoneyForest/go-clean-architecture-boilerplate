package testhelper

import (
	"context"
	"database/sql"

	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/aws"
	mysqlgw "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql"
	redisgw "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/redis"
	sqsgw "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs"
	"github.com/redis/go-redis/v9"
)

type Gateway struct {
	MySQLClient *sql.DB
	RedisClient *redis.Client
	SQSClient   *sqsgw.SQSClient
}

type TestEnvironment struct {
	Environment        string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBDatabase         string
	RedisHost          string
	RedisPort          string
	RedisPassword      string
	AWSRegion          string
	AWSEndpoint        string
	SQSQueueNameSample string
}

func Setup(ctx context.Context) (*Gateway, error) {
	e := &TestEnvironment{
		Environment:        "test",
		DBHost:             "localhost",
		DBPort:             "3306",
		DBUser:             "root",
		DBPassword:         "password",
		DBDatabase:         "maindb",
		RedisHost:          "localhost",
		RedisPort:          "6379",
		RedisPassword:      "password",
		AWSRegion:          "ap-northeast-1",
		AWSEndpoint:        "http://localhost:4566",
		SQSQueueNameSample: "sample_queue",
	}

	mysqlClient, err := mysqlgw.InitDB(ctx, mysqlgw.DBConfig{
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
	redisClient, err := redisgw.InitRedis(ctx, redisgw.RedisConfig{
		Environment: e.Environment,
		Host:        e.RedisHost,
		Port:        e.RedisPort,
		Password:    e.RedisPassword,
	})
	if err != nil {
		return nil, err
	}
	sqsClient, err := sqsgw.InitSQS(ctx, aws.AWSConfig{
		Environment: e.Environment,
		Region:      e.AWSRegion,
		Endpoint:    e.AWSEndpoint,
	}, sqsgw.SQSConfig{
		QueueNames: map[sqsgw.Key]string{sqsgw.SQSKeySample: e.SQSQueueNameSample},
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
