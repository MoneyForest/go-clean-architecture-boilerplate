package model

import (
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        uuid.UUID `validate:"required"`
	Email     string    `validate:"required,email"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
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
	validate := validator.New()
	if err := validate.Struct(u); err != nil {
		return err
	}
	return nil
}
