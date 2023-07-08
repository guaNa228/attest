-- name: GetPreAttestationData :many
SELECT u.id as student_id,
    sa.id as semester_activity_id
FROM users u,
    semester_activity sa
WHERE u.group_id = sa.group_id;
-- name: GetPreAttestationData :many
SELECT u.id as student,
    w.id as workload
FROM users u,
    workloads w
WHERE u.group_id = w.group_id;
-- name: GetAttestationData :many
SELECT a.id,
    u.name student,
    s.code,
    s.name,
    g.subcode,
    c.name class,
    a.result,
    a.month,
    a.comment
FROM workloads w,
    attestation a,
    users u,
    groups g,
    classes c,
    streams s
WHERE w.teacher = $1
    and a.workload = w.id
    and a.student = u.id
    and g.id = w.group_id
    and g.stream = s.id
    and c.id = w.class;
-- name: UpdateAttestationRow :exec
UPDATE attestation
SET result = $2,
    comment = $3
WHERE id = $1;
-- name: GetPreAttestationData :many
SELECT u.id as student,
    w.id as workload
FROM users u,
    workloads w
WHERE u.group_id = w.group_id;
-- name: GetStudentsAttestationData :many
SELECT c.name class,
    a.result,
    a.month,
    a.comment
FROM attestation a,
    worloads w,
    classes c
WHERE a.student = $1
    and a.workload = w.id
    and w.class = c.id;