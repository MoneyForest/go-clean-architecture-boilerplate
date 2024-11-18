package input

import "github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"

type CreateUserInput struct {
	Email string `json:"email"`
}

type GetUserInput struct {
	ID uuid.UUID `json:"id"`
}

type ListUserInput struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type UpdateUserInput struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type DeleteUserInput struct {
	ID uuid.UUID `json:"id"`
}

type ProcessMessageInput struct {
	QueueURL string `json:"queue_url"`
}
