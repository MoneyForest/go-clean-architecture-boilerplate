package response

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
