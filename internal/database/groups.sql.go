// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: query.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createGroup = `-- name: CreateGroup :one
INSERT INTO groups(
        id,
        created_at,
        updated_at,
        name,
        code
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
RETURNING id, created_at, updated_at, name, code
`

type CreateGroupParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Code      string
}

func (q *Queries) CreateGroup(ctx context.Context, arg CreateGroupParams) (Group, error) {
	row := q.db.QueryRowContext(ctx, createGroup,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Code,
	)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Code,
	)
	return i, err
}

const deleteGroupByCode = `-- name: DeleteGroupByCode :exec
DELETE FROM groups
WHERE code = $1
`

func (q *Queries) DeleteGroupByCode(ctx context.Context, code string) error {
	_, err := q.db.ExecContext(ctx, deleteGroupByCode, code)
	return err
}

const deleteGroupByID = `-- name: DeleteGroupByID :exec
DELETE FROM groups
WHERE id = $1
`

func (q *Queries) DeleteGroupByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteGroupByID, id)
	return err
}