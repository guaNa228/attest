-- name: CreateGroup :one
INSERT INTO groups(
        id,
        created_at,
        updated_at,
        name,
        code
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
RETURNING *;
-- name: DeleteGroupByID :one
DELETE FROM groups
WHERE id = $1;
-- name: DeleteGroupByCode :one
DELETE FROM groups
WHERE code = $1;