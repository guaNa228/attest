-- name: CreateSemesterActivity :one
INSERT INTO semester_activity(
        id,
        created_at,
        updated_at,
        group_id,
        class_id,
        teacher_id
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: UpdateSemesterActivityById :one
UPDATE semester_activity
SET group_id = $2,
    class_id = $3,
    teacher_id = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
-- name: DeleteSemesterActivityById :exec
DELETE FROM semester_activity
WHERE code = $1;