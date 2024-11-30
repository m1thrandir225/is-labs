-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    secret_key,
    counter
) VALUES (
        ?,
        ?,
        ?,
        ?
) RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?;

-- name: UpdateUserCounter :exec
UPDATE users
SET counter = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
