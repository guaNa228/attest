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
-- name: GetGroupByFullCode :one
SELECT g.id
FROM groups g,
    streams s
WHERE g.stream = s.id
    and s.code = $1
    and g.subcode = $2;
-- name: GetFileData :many 
SELECT a.code,
    a.stream,
    a.name,
    a.email
from (
        SELECT s.code || '/' || g.subcode as code,
            s.name as stream,
            u.name,
            u.email
        from users u,
            groups g,
            streams s
        where u.role = 'student'
            and u.group_id = g.id
            and g.stream = s.id
    ) a
group by a.stream,
    a.code,
    a.name,
    a.email