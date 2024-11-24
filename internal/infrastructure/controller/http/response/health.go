package response

type HealthResponse struct {
	Status string `json:"status"`
}

type DeepHealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
