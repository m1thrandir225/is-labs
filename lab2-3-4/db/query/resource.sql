-- name: CreateResource :one
INSERT INTO resources (name, org_id)
VALUES (?, ?)
RETURNING *;

-- name: GetResource :one
SELECT *
FROM resources
WHERE id = ?
LIMIT 1;

-- name: ListOrganizationResources :many
SELECT *
FROM resources
WHERE org_id = ?
ORDER BY created_at;