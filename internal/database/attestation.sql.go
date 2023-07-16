// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: query.sql

package db

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

const getPreAttestationData = `-- name: GetPreAttestationData :many
SELECT u.id as student,
    w.id as workload
FROM users u,
    workloads w
WHERE u.group_id = w.group_id
`

type GetPreAttestationDataRow struct {
	Student  uuid.UUID
	Workload uuid.UUID
}

func (q *Queries) GetPreAttestationData(ctx context.Context) ([]GetPreAttestationDataRow, error) {
	rows, err := q.db.QueryContext(ctx, getPreAttestationData)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPreAttestationDataRow
	for rows.Next() {
		var i GetPreAttestationDataRow
		if err := rows.Scan(&i.Student, &i.Workload); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAttestationRow = `-- name: UpdateAttestationRow :exec
UPDATE attestation
SET result = $2,
    comment = $3
WHERE id = $1
`

type UpdateAttestationRowParams struct {
	ID      uuid.UUID
	Result  sql.NullInt32
	Comment sql.NullString
}

func (q *Queries) UpdateAttestationRow(ctx context.Context, arg UpdateAttestationRowParams) error {
	_, err := q.db.ExecContext(ctx, updateAttestationRow, arg.ID, arg.Result, arg.Comment)
	return err
}

const getStudentsAttestationData = `-- name: GetStudentsAttestationData :many
SELECT c.name class,
    a.result,
    a.month,
    a.comment
FROM attestation a,
    workloads w,
    classes c
WHERE a.student = $1
    and a.workload = w.id
    and w.class = c.id
`

type GetStudentsAttestationDataRow struct {
	Class   string `json:"class"`
	Result  sql.NullInt32 `json:"result"`
	Month   MonthEnum `json:"month"`
	Comment sql.NullString `json:"comment"`
}

func (q *Queries) GetStudentsAttestationData(ctx context.Context, student uuid.UUID) ([]GetStudentsAttestationDataRow, error) {
	rows, err := q.db.QueryContext(ctx, getStudentsAttestationData, student)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetStudentsAttestationDataRow
	for rows.Next() {
		var i GetStudentsAttestationDataRow
		if err := rows.Scan(
			&i.Class,
			&i.Result,
			&i.Month,
			&i.Comment,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const clearAttestation = `-- name: ClearAttestation :exec
DELETE from attestation
where month = $1
`

func (q *Queries) ClearAttestation(ctx context.Context, month MonthEnum) error {
	_, err := q.db.ExecContext(ctx, clearAttestation, month)
	return err
}

const getTeachersAttestationData = `-- name: GetTeachersAttestationData :many
SELECT w.id,
	s.code || '/' || g.subcode as "group_code",
	s.name,
    c.name class
FROM  workloads w,
    groups g,
    classes c,
	streams s
WHERE w.teacher = $1
    and g.id = w.group_id
	and g.stream = s.id
    and c.id = w.class
	and (SELECT COUNT(*)
		from users
		where group_id=g.id)>0
`

type GetTeachersAttestationDataRow struct {
	ID        uuid.UUID `json:"id"`
	GroupCode string `json:"group_code"`
	Name      string `json:"stream"`
	Class     string `json:"class"`
}

func (q *Queries) GetTeachersAttestationData(ctx context.Context, teacher uuid.UUID) ([]GetTeachersAttestationDataRow, error) {
	rows, err := q.db.QueryContext(ctx, getTeachersAttestationData, teacher)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTeachersAttestationDataRow
	for rows.Next() {
		var i GetTeachersAttestationDataRow
		if err := rows.Scan(
			&i.ID,
			&i.GroupCode,
			&i.Name,
			&i.Class,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWorkloadAttestationData = `-- name: GetWorkloadAttestationData :many
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
    and u.id = a.student
`

type GetWorkloadAttestationDataRow struct {
	ID      uuid.UUID `json:"id"`
	Name    string `json:"student"`
	Month   MonthEnum `json:"month"`
	Result  sql.NullInt32 `json:"result"`
	Comment sql.NullString `json:"comment"`
}

func (q *Queries) GetWorkloadAttestationData(ctx context.Context, id uuid.UUID) ([]GetWorkloadAttestationDataRow, error) {
	rows, err := q.db.QueryContext(ctx, getWorkloadAttestationData, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWorkloadAttestationDataRow
	for rows.Next() {
		var i GetWorkloadAttestationDataRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Month,
			&i.Result,
			&i.Comment,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}