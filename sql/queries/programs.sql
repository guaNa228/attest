-- name: CreateProgram :one
INSERT INTO programs (id, name, max_courses)
VALUES ($1, $2, $3)
RETURNING *;
-- name: GetProgramsNumber :one
SELECT COUNT(*)
FROM programs;