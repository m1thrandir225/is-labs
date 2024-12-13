// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    otp_secret,
    is_2fa_enabled
) VALUES (
        ?,
        ?,
        ?,
        ?
) RETURNING id, email, password_hash, otp_secret, is_2fa_enabled
`

type CreateUserParams struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	OtpSecret    string `json:"otp_secret"`
	Is2faEnabled bool   `json:"is_2fa_enabled"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Email,
		arg.PasswordHash,
		arg.OtpSecret,
		arg.Is2faEnabled,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.OtpSecret,
		&i.Is2faEnabled,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password_hash, otp_secret, is_2fa_enabled
FROM users
WHERE email = ?
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.OtpSecret,
		&i.Is2faEnabled,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, email, password_hash, otp_secret, is_2fa_enabled
FROM users
WHERE id = ?
`

func (q *Queries) GetUserById(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.OtpSecret,
		&i.Is2faEnabled,
	)
	return i, err
}
