-- name: CreateGroup :one
INSERT INTO groups(
        id,
        created_at,
        updated_at,
        subcode,
        stream,
        course,
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
RETURNING *;
-- name: DeleteGroupByID :exec
DELETE FROM groups
WHERE id = $1;
-- name: DeleteGroupByCode :exec
DELETE FROM groups
WHERE code = $1;