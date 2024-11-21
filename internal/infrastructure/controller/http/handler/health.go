package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/http/dto"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/http/dto/response"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type DeepHealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Details   map[string]ServiceInfo `json:"details"`
}

type ServiceInfo struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, HealthResponse{
		Status: "ok",
	})
}

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

	response.WriteJSON(w, statusCode, dto.DeepHealthResponse{
		Status: status,
	})
}
