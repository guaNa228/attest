-- name: CreateUser :one
INSERT INTO users(
        id,
        created_at,
        updated_at,
        name,
        login,
        password,
        role,
        group_id,
        teacher_id
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9
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
-- name: UpdateStudentsGroup :one
UPDATE users
SET group_id = $2,
    updated_at = NOW()
WHERE id = $1
    and role = "student"
RETURNING *;
-- name: IfLoginDuplicates :one
select exists(
        select 1
        from users
        where login ~ $1
    );
-- name: NumberOfDuplicatedUsers :one
select COUNT(*)
from users
where login ~ $1;
-- name: GetTeacherIDByNameAndTeacherId :one
select id
from users
where name = $1
    and teacher_id = $2;
-- name: GetUserByEmail :one
select id
from users
where email = $1;
-- name: DeleteSemesterUsers :exec
DELETE from users
WHERE email is null
    and role = 'student';
-- name: RemoveGroupID :exec
UPDATE users
SET group_id = NULL
WHERE group_id IS NOT NULL
    and role = 'student';
-- name: GetTeachersWithUniqueName :many
SELECT id,
    name
FROM users
WHERE role = 'teacher'
    AND name IN (
        SELECT name
        FROM users
        WHERE role = 'teacher'
        GROUP BY name
        HAVING COUNT(*) = 1
    )
    and email is null;
-- name: GetFullUserByEmail :one
select *
from users
where email = $1;