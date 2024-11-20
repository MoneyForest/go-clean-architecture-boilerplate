package entity

import (
	"time"
)

type MatchingEntity struct {
	ID        string    `db:"id"`
	MeID      string    `db:"me_id"`
	PartnerID string    `db:"partner_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
