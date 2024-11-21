package entity

import (
	"time"
)

type UserEntity struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
