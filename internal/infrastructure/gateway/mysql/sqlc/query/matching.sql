-- name: GetMatching :one
SELECT * FROM `matching`
WHERE id = ? LIMIT 1;

-- name: GetMatchingByParticipants :one
SELECT * FROM `matching`
WHERE me_id = ? AND partner_id = ?
LIMIT 1;

-- name: ListMatchingsByUser :many
SELECT * FROM `matching`
WHERE me_id = ? OR partner_id = ?
LIMIT ? OFFSET ?;

-- name: ExistsMatching :one
SELECT EXISTS(
    SELECT 1 FROM `matching` WHERE id = ?
);

-- name: CreateMatching :execresult
INSERT INTO `matching` (
    id,
    me_id,
    partner_id,
    `status`,
    created_at,
    updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?
);

-- name: UpdateMatching :execresult
UPDATE `matching`
SET
    `status` = ?,
    updated_at = ?
WHERE id = ?;

-- name: DeleteMatching :exec
DELETE FROM `matching`
WHERE id = ?;
