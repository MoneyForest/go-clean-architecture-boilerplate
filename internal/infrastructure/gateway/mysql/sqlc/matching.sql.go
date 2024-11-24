// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: matching.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const CreateMatching = `-- name: CreateMatching :execresult
INSERT INTO ` + "`" + `matching` + "`" + ` (
    id,
    me_id,
    partner_id,
    ` + "`" + `status` + "`" + `,
    created_at,
    updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?
)
`

type CreateMatchingParams struct {
	ID        string    `json:"id"`
	MeID      string    `json:"me_id"`
	PartnerID string    `json:"partner_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreateMatching(ctx context.Context, arg CreateMatchingParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateMatching,
		arg.ID,
		arg.MeID,
		arg.PartnerID,
		arg.Status,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
}

const DeleteMatching = `-- name: DeleteMatching :exec
DELETE FROM ` + "`" + `matching` + "`" + `
WHERE id = ?
`

func (q *Queries) DeleteMatching(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, DeleteMatching, id)
	return err
}

const ExistsMatching = `-- name: ExistsMatching :one
SELECT EXISTS(
    SELECT 1 FROM ` + "`" + `matching` + "`" + ` WHERE id = ?
)
`

func (q *Queries) ExistsMatching(ctx context.Context, id string) (bool, error) {
	row := q.db.QueryRowContext(ctx, ExistsMatching, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const GetMatching = `-- name: GetMatching :one
SELECT id, me_id, partner_id, status, created_at, updated_at FROM ` + "`" + `matching` + "`" + `
WHERE id = ? LIMIT 1
`

func (q *Queries) GetMatching(ctx context.Context, id string) (Matching, error) {
	row := q.db.QueryRowContext(ctx, GetMatching, id)
	var i Matching
	err := row.Scan(
		&i.ID,
		&i.MeID,
		&i.PartnerID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetMatchingByParticipants = `-- name: GetMatchingByParticipants :one
SELECT id, me_id, partner_id, status, created_at, updated_at FROM ` + "`" + `matching` + "`" + `
WHERE me_id = ? AND partner_id = ?
LIMIT 1
`

type GetMatchingByParticipantsParams struct {
	MeID      string `json:"me_id"`
	PartnerID string `json:"partner_id"`
}

func (q *Queries) GetMatchingByParticipants(ctx context.Context, arg GetMatchingByParticipantsParams) (Matching, error) {
	row := q.db.QueryRowContext(ctx, GetMatchingByParticipants, arg.MeID, arg.PartnerID)
	var i Matching
	err := row.Scan(
		&i.ID,
		&i.MeID,
		&i.PartnerID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const ListMatchingsByUser = `-- name: ListMatchingsByUser :many
SELECT id, me_id, partner_id, status, created_at, updated_at FROM ` + "`" + `matching` + "`" + `
WHERE me_id = ? OR partner_id = ?
LIMIT ? OFFSET ?
`

type ListMatchingsByUserParams struct {
	MeID      string `json:"me_id"`
	PartnerID string `json:"partner_id"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

func (q *Queries) ListMatchingsByUser(ctx context.Context, arg ListMatchingsByUserParams) ([]Matching, error) {
	rows, err := q.db.QueryContext(ctx, ListMatchingsByUser,
		arg.MeID,
		arg.PartnerID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Matching{}
	for rows.Next() {
		var i Matching
		if err := rows.Scan(
			&i.ID,
			&i.MeID,
			&i.PartnerID,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const UpdateMatching = `-- name: UpdateMatching :execresult
UPDATE ` + "`" + `matching` + "`" + `
SET
    ` + "`" + `status` + "`" + ` = ?,
    updated_at = ?
WHERE id = ?
`

type UpdateMatchingParams struct {
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        string    `json:"id"`
}

func (q *Queries) UpdateMatching(ctx context.Context, arg UpdateMatchingParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, UpdateMatching, arg.Status, arg.UpdatedAt, arg.ID)
}
