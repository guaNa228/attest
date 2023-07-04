-- name: CreateClass :one
INSERT INTO classes(id, name)
VALUES ($1, $2)
RETURNING *;
-- name: DeleteClassByID :exec
DELETE FROM classes
WHERE id = $1;
-- name: ClearClassesTable :exec
DELETE from classes;