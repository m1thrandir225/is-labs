-- name: AddUserToOrganization :one
INSERT INTO user_organizations (user_id, org_id, role_id)
VALUES (?,
        ?,
        ?)
RETURNING *;

-- name: GetUserOrganization :one
SELECT *
FROM user_organizations
WHERE user_id = ?
  AND org_id = ?;

-- name: ListUserOrganizations :many
SELECT *
FROM user_organizations
WHERE user_id = ?
ORDER BY created_at;

-- name: RemoveUserFromOrganization :exec
DELETE FROM user_organizations
WHERE org_id = ? AND user_id = ?;