// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users(
    id, username, email, password
) VALUES (?, ?, ?, ?)
`

type CreateUserParams struct {
	ID       string
	Username string
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser,
		arg.ID,
		arg.Username,
		arg.Email,
		arg.Password,
	)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserById = `-- name: GetUserById :one
SELECT id, username, email, password, created_at FROM users
WHERE id = ? LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
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
LIMIT ? OFFSET ?
`

type ListUsersParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
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

const retrieveIdByEmail = `-- name: RetrieveIdByEmail :one
SELECT id FROM users
WHERE email = ? LIMIT 1
`

func (q *Queries) RetrieveIdByEmail(ctx context.Context, email string) (string, error) {
	row := q.db.QueryRowContext(ctx, retrieveIdByEmail, email)
	var id string
	err := row.Scan(&id)
	return id, err
}

const retrieveIdByUsername = `-- name: RetrieveIdByUsername :one
SELECT id FROM users
WHERE username = ? LIMIT 1
`

func (q *Queries) RetrieveIdByUsername(ctx context.Context, username string) (string, error) {
	row := q.db.QueryRowContext(ctx, retrieveIdByUsername, username)
	var id string
	err := row.Scan(&id)
	return id, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users SET username = ?, email = ?, password = ?
WHERE id=?
`

type UpdateUserParams struct {
	Username string
	Email    string
	Password string
	ID       string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.Username,
		arg.Email,
		arg.Password,
		arg.ID,
	)
	return err
}
