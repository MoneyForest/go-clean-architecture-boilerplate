package interactor

import (
	"context"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/repository"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/port"
)

type HealthInteractor struct {
	healthRepo  repository.HealthRepository
	healthCache repository.HealthCacheRepository
}

func NewHealthInteractor(healthRepo repository.HealthRepository, healthCache repository.HealthCacheRepository) HealthInteractor {
	return HealthInteractor{
		healthRepo:  healthRepo,
		healthCache: healthCache,
	}
}

func (i HealthInteractor) Check(ctx context.Context, input *port.CheckHealthInput) (*port.CheckHealthOutput, error) {
	return &port.CheckHealthOutput{
		Status: "ok",
	}, nil
}

func (i HealthInteractor) DeepCheck(ctx context.Context, input *port.DeepCheckHealthInput) (*port.DeepCheckHealthOutput, error) {
	status := "ok"
	message := []string{}

	if err := i.healthRepo.Ping(ctx); err != nil {
		message = append(message, err.Error())
		status = "error"
	}
	if err := i.healthCache.Ping(ctx); err != nil {
		message = append(message, err.Error())
		status = "error"
	}

	return &port.DeepCheckHealthOutput{
		Status:  status,
		Message: message,
	}, nil
}
