package entity

import (
	"time"
)

type MatchEntity struct {
	ID        string    `db:"id"`
	MeID      string    `db:"user1_id"`
	PartnerID string    `db:"user2_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
