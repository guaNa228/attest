-- name: GetPreAttestationData :many
SELECT u.id as student_id,
    sa.id as semester_activity_id
FROM users u,
    semester_activity sa
WHERE u.group_id = sa.group_id;
-- name: GetAttestationData :many
SELECT a.id,
    u.name student,
    g.code,
    c.name class,
    a.result,
    a.month,
    a.comment
FROM semester_activity sa,
    attestation a,
    users u,
    groups g,
    classes c
WHERE sa.teacher_id = $1
    and a.student_id = u.id
    and g.id = sa.group_id
    and c.id = sa.class_id;