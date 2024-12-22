-- name: CreateResource :one
INSERT INTO resources (name, org_id)
VALUES (?, ?)
RETURNING *;

-- name: GetResource :one
SELECT *
FROM resources
WHERE id = ?
LIMIT 1;

-- name: UpdateResource :one
UPDATE resources
SET
    name = ?
WHERE id = ?
RETURNING *;

-- name: DeleteResource :exec
DELETE FROM resources
WHERE id = ?;

-- name: ListOrganizationResources :many
SELECT *
FROM resources
WHERE org_id = ?
ORDER BY created_at;