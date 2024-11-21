package model

import (
	"errors"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

var (
	ErrUserEmailIsRequired = errors.New("user email is required")
)

type User struct {
	ID        uuid.UUID
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type InputUserParams struct {
	ID    uuid.UUID
	Email string
}

func NewUser(params InputUserParams) *User {
	if params.ID == uuid.Nil() {
		params.ID = uuid.New()
	}
	return &User{
		ID:        params.ID,
		Email:     params.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) Validate() error {
	if u.Email == "" {
		return ErrUserEmailIsRequired
	}
	return nil
}
