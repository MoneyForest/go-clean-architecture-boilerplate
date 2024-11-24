package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/http/marshaller"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/http/response"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/interactor"
)

// @title			Health Handler
// @description	Handles HTTP requests for system health checks
type HealthHandler struct {
	HealthInteractor interactor.HealthInteractor
}

// @Summary		Get system health status
// @Description	Performs a basic health check of the system
// @Tags			health
// @Accept			json
// @Produce		json
// @Success		200	{object}	response.HealthResponse	"System is healthy"
// @Failure		503	{object}	response.HealthResponse	"Service unavailable"
// @Router			/check [get]
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	output, err := h.HealthInteractor.Check(ctx, marshaller.ToCheckHealthInput())
	if err != nil {
		response.WriteJSON(w, http.StatusServiceUnavailable, response.HealthResponse{
			Status: "error",
		})
	}
	response.WriteJSON(w, http.StatusOK, marshaller.ToCheckHealthResponse(output))
}

// @Summary		Get detailed system health status
// @Description	Performs a deep health check including dependent services
// @Tags			health
// @Accept			json
// @Produce		json
// @Success		200	{object}	response.DeepHealthResponse	"System is healthy"
// @Failure		503	{object}	response.DeepHealthResponse	"System is unhealthy"
// @Router			/deep_check [get]
func (h *HealthHandler) DeepCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	output, err := h.HealthInteractor.DeepCheck(ctx, marshaller.ToDeepCheckHealthInput())
	if err != nil {
		response.WriteJSON(w, http.StatusServiceUnavailable, response.DeepHealthResponse{
			Status:  output.Status,
			Message: strings.Join(output.Message, "\n"),
		})
	}
	response.WriteJSON(w, http.StatusOK, marshaller.ToDeepCheckHealthResponse(output))
}
