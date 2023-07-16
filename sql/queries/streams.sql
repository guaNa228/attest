-- name: ClearStreamsTable :exec
DELETE from streams;
-- name: GetAllStreams :many
SELECT DISTINCT s.id,
    s.name
from attestation a,
    workloads w,
    groups g,
    streams s
where w.group_id = g.id
    and g.stream = s.id
    and a.workload = w.id;
-- name: GetStreamByID :one
SELECT name
from streams
where id = $1;