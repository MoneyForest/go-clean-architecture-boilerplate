package dependency

import (
	"context"
	"fmt"

	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/environment"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql"
	mysqlRepo "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql/repository"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/redis"
	redisRepo "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/redis/repository"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs"
	sqsRepo "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs/repository"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/interactor"
	"github.com/caarlos0/env/v10"
)

type Dependency struct {
	Environment    *environment.Environment
	UserInteractor interactor.UserInteractor
}

func Inject(ctx context.Context) (*Dependency, error) {
	e := &environment.Environment{}
	if err := env.Parse(e); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	mysqlClient, err := mysql.InitDB(ctx, mysql.DBConfig{
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
	redisClient, err := redis.InitRedis(ctx, redis.RedisConfig{
		Environment: e.Environment,
		Host:        e.RedisHost,
		Port:        e.RedisPort,
		Password:    e.RedisPassword,
	})
	if err != nil {
		return nil, err
	}
	sqsClient, err := sqs.InitSQS(ctx, sqs.SQSConfig{
		Environment: e.Environment,
		Region:      e.AWSRegion,
		Endpoint:    e.SQSEndpoint,
	})
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	mysqlUserRepository := mysqlRepo.NewUserMySQLRepository(mysqlClient)
	redisUserRepository := redisRepo.NewUserRedisRepository(redisClient)
	sqsRepository := sqsRepo.NewSQSRepository(sqsClient)

	// Initialize interactor
	userInteractor := interactor.NewUserInteractor(mysqlUserRepository, redisUserRepository, sqsRepository)

	return &Dependency{
		Environment:    e,
		UserInteractor: userInteractor,
	}, nil
}
