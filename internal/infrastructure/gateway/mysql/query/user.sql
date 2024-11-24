-- name: GetUser :one
SELECT * FROM `user`
WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM `user`
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateUser :execresult
INSERT INTO `user` (
    id,
    email,
    created_at,
    updated_at
) VALUES (
    ?, ?, ?, ?
);

-- name: UpdateUser :execresult
UPDATE `user`
SET
    email = ?,
    updated_at = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM `user`
WHERE id = ?;

-- name: ExistsUser :one
SELECT EXISTS(
    SELECT 1 FROM `user` WHERE id = ?
);

-- name: CountUsers :one
SELECT COUNT(*) FROM `user`;
