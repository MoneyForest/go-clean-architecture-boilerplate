package marshaller

import (
	"strings"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/http/response"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/port"
)

// Input Marshalling
func ToCheckHealthInput() *port.CheckHealthInput {
	return &port.CheckHealthInput{}
}

func ToDeepCheckHealthInput() *port.DeepCheckHealthInput {
	return &port.DeepCheckHealthInput{}
}

// Output Marshalling
func ToCheckHealthResponse(output *port.CheckHealthOutput) response.HealthResponse {
	return response.HealthResponse{
		Status: output.Status,
	}
}

func ToDeepCheckHealthResponse(output *port.DeepCheckHealthOutput) response.DeepHealthResponse {
	return response.DeepHealthResponse{
		Status:  output.Status,
		Message: strings.Join(output.Message, "\n"),
	}
}
