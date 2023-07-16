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
-- name: GetTeachersAttestationData :many
SELECT w.id,
    s.code || '/' || g.subcode as "group_code",
    s.name,
    c.name class
FROM workloads w,
    groups g,
    classes c,
    streams s
WHERE w.teacher = $1
    and g.id = w.group_id
    and g.stream = s.id
    and c.id = w.class
    and (
        SELECT COUNT(*)
        from users
        where group_id = g.id
    ) > 0;
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
-- name: ClearAttestation :exec
DELETE from attestation
where month = $1;
-- name: GetWorkloadAttestationData :many
SELECT a.id,
    u.name,
    a.month,
    a.result,
    a.comment
from attestation a,
    workloads w,
    users u
WHERE w.id = $1
    and a.workload = w.id
    and u.id = a.student;
-- name: GetUnattestedStudents :many
SELECT a.student id,
    COUNT(*) as "cnt"
from attestation a,
    workloads w,
    groups g
where a.workload = w.id
    and w.group_id = g.id
    and g.stream = $1
    and a.result < $2
    and a.month = $3
group by a.student;
-- name: GetUnderachieversData :many
SELECT u.name student,
    c.name class,
    attestation.result res,
    s.code || '/' || g.subcode group_code
from (
        SELECT student id
        from attestation a,
            workloads w,
            groups g
        where a.workload = w.id
            and w.group_id = g.id
            and g.stream = $1
            and a.result < $2
            and a.month = $4
        group by a.student
        HAVING COUNT(*) >= $3
    ) st,
    workloads w,
    users u,
    classes c,
    groups g,
    streams s,
    attestation
where attestation.student = st.id
    and u.id = st.id
    and attestation.workload = w.id
    and w.class = c.id
    and w.group_id = g.id
    and g.stream = s.id
    and attestation.result < $2;
-- name: GetOpenedMonths :many
SELECT DISTINCT month
from attestation;