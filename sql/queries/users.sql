-- name: CreateUser :one
INSERT INTO users(
        id,
        created_at,
        updated_at,
        name,
        login,
        password,
        role
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7
    )
RETURNING *;
-- name: GetUserByCredentials :one
SELECT *
FROM users
WHERE login = $1
    and password = $2;
-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;