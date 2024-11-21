package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/http/response"
)

// @title			Health Handler
// @description	Handles HTTP requests for health checks
type HealthHandler struct {
	db *sql.DB
}

// @Summary	Get health status
// @Tags		health
// @Accept		json
// @Produce	json
// @Success	200	{object}	response.HealthResponse
// @Router		/health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, response.HealthResponse{
		Status: "ok",
	})
}

// @Summary	Get deep health status
// @Tags		health
// @Accept		json
// @Produce	json
// @Success	200	{object}	response.DeepHealthResponse
// @Failure	503	{object}	error.DomainError
// @Router		/health/deep [get]
func (h *HealthHandler) Deep(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	status := "ok"
	if err := h.db.PingContext(ctx); err != nil {
		status = "error"
	}

	statusCode := http.StatusOK
	if status != "ok" {
		statusCode = http.StatusServiceUnavailable
	}

	response.WriteJSON(w, statusCode, response.DeepHealthResponse{
		Status: status,
	})
}
