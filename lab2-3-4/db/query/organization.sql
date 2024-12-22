-- name: CreateOrganization :one
INSERT INTO organizations (name)
VALUES (?)
RETURNING *;

-- name: GetOrganization :one
SELECT *
FROM organizations
WHERE id = ?;

-- name: DeleteOrganization :exec
DELETE
FROM organizations
WHERe id = ?;