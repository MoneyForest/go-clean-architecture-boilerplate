package port

type CheckHealthInput struct{}

type CheckHealthOutput struct {
	Status string `json:"status"`
}

type DeepCheckHealthInput struct{}

type DeepCheckHealthOutput struct {
	Status  string   `json:"status"`
	Message []string `json:"message"`
}
