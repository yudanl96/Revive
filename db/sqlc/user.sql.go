// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :execresult
INSERT INTO users(
    username, email, password
) VALUES (?, ?, ?)
`

type CreateUserParams struct {
	Username string
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createUser, arg.Username, arg.Email, arg.Password)
}

const deleteUser = `-- name: DeleteUser :execresult
DELETE FROM users WHERE id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id string) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteUser, id)
}

const getUser = `-- name: GetUser :one
SELECT id, username, email, password, created_at FROM users
WHERE id = ? LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, username, email, password, created_at FROM users 
ORDER BY id
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.CreatedAt,
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

const updateUser = `-- name: UpdateUser :execresult
UPDATE users SET username = ?, email = ?, password = ? 
WHERE id=?
`

type UpdateUserParams struct {
	Username string
	Email    string
	Password string
	ID       string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUser,
		arg.Username,
		arg.Email,
		arg.Password,
		arg.ID,
	)
}
