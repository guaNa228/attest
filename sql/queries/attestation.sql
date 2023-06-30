-- name: GetPreAttestationData :many
SELECT u.id as student_id,
    sa.id as semester_activity_id
FROM users u,
    semester_activity sa
WHERE u.group_id = sa.group_id;