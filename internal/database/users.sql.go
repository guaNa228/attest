// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: query.sql

package db

import (
	"context"
	"time"
	"database/sql"
	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
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
RETURNING id, created_at, updated_at, name, login, password, role, teacher_id, group_id
`

type CreateUserParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Login     string
	Password  string
	Role      string
	GroupID   uuid.NullUUID
	TeacherID sql.NullInt32
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Login,
		arg.Password,
		arg.Role,
		arg.GroupID,
		arg.TeacherID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Login,
		&i.Password,
		&i.Role,
		&i.TeacherID,
		&i.GroupID,
	)
	return i, err
}

const getUserByCredentials = `-- name: GetUserByCredentials :one
SELECT id, created_at, updated_at, name, login, password, role, group_id
FROM users
WHERE login = $1
    and password = $2
`

type GetUserByCredentialsParams struct {
	Login    string
	Password string
}

func (q *Queries) GetUserByCredentials(ctx context.Context, arg GetUserByCredentialsParams) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByCredentials, arg.Login, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Login,
		&i.Password,
		&i.Role,
		&i.GroupID,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, created_at, updated_at, name, login, password, role, group_id
FROM users
WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Login,
		&i.Password,
		&i.Role,
		&i.GroupID,
	)
	return i, err
}

const updateStudentsGroup = `-- name: UpdateStudentsGroup :one
UPDATE users
SET group_id = $2,
    updated_at = NOW()
WHERE id = $1
    and role = "student"
RETURNING id, created_at, updated_at, name, login, password, role, group_id
`

type UpdateStudentsGroupParams struct {
	ID      uuid.UUID
	GroupID uuid.NullUUID
}

func (q *Queries) UpdateStudentsGroup(ctx context.Context, arg UpdateStudentsGroupParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateStudentsGroup, arg.ID, arg.GroupID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Login,
		&i.Password,
		&i.Role,
		&i.GroupID,
	)
	return i, err
}

const ifLoginDuplicates = `-- name: IfLoginDuplicates :one
select exists(
        select 1
        from users
        where login ~ $1
    )
`

func (q *Queries) IfLoginDuplicates(ctx context.Context, login string) (bool, error) {
	row := q.db.QueryRowContext(ctx, ifLoginDuplicates, login)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const numberOfDuplicatedUsers = `-- name: NumberOfDuplicatedUsers :one
select COUNT(*)
from users
where login ~ $1
`

func (q *Queries) NumberOfDuplicatedUsers(ctx context.Context, login string) (int64, error) {
	row := q.db.QueryRowContext(ctx, numberOfDuplicatedUsers, login)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getTeacherIDByNameAndTeacherId = `-- name: GetTeacherIDByNameAndTeacherId :one
select id
from users
where name=$1 and teacher_id=$2
`

type GetTeacherIDByNameAndTeacherIdParams struct {
	Name      string
	TeacherID sql.NullInt32
}

func (q *Queries) GetTeacherIDByNameAndTeacherId(ctx context.Context, arg GetTeacherIDByNameAndTeacherIdParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getTeacherIDByNameAndTeacherId, arg.Name, arg.TeacherID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
select id
from users
where email=$1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email sql.NullString) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}


const deleteSemesterUsers = `-- name: DeleteSemesterUsers :exec
DELETE from users
WHERE email is null
    and role = 'student'
`

func (q *Queries) DeleteSemesterUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteSemesterUsers)
	return err
}

const removeGroupID = `-- name: RemoveGroupID :exec
UPDATE users
SET group_id = NULL
WHERE group_id IS NOT NULL
    and role = 'student'
`

func (q *Queries) RemoveGroupID(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, removeGroupID)
	return err
}

const getTeachersWithUniqueName = `-- name: GetTeachersWithUniqueName :many
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
    and email is null
`

type GetTeachersWithUniqueNameRow struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) GetTeachersWithUniqueName(ctx context.Context) ([]*GetTeachersWithUniqueNameRow, error) {
	rows, err := q.db.QueryContext(ctx, getTeachersWithUniqueName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetTeachersWithUniqueNameRow
	for rows.Next() {
		var i GetTeachersWithUniqueNameRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFullUserByEmail = `-- name: GetFullUserByEmail :one
select id, created_at, updated_at, name, login, password, role, teacher_id, group_id, email
from users
where email = $1
`

func (q *Queries) GetFullUserByEmail(ctx context.Context, email sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getFullUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Login,
		&i.Password,
		&i.Role,
		&i.TeacherID,
		&i.GroupID,
		&i.Email,
	)
	return i, err
}
