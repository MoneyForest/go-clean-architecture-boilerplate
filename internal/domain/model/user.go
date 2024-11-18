package model

import (
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email string) *User {
	return &User{
		ID:        uuid.New(),
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
