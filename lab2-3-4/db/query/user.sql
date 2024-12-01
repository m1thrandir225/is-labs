-- name: CreateUser :one
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
) RETURNING *;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
