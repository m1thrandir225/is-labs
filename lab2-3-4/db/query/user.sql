-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    otp_secret,
    is_2fa_enabled,
    role
) VALUES (
        ?,
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

-- name: UpdateUserRole :exec
UPDATE users SET role = ? WHERE id = ?;

-- name: GetUserRole :one
SELECT role FROM users WHERE id = ?;
