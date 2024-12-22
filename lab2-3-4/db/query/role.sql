-- name: CreateRole :one
INSERT INTO roles (name, org_id)
VALUES (?, ?)
RETURNING *;

-- name: GetRole :one
SELECT *
FROM roles
WHERE id = ?
LIMIT 1;

-- name: ListOrganizationRoles :many
SELECT *
FROM roles
WHERE org_id = ?
ORDER BY created_at;

-- name: UpdateRole :one
UPDATE roles
SET name = ?
WHERE id = ?
RETURNING *;


-- name: GetUserRole :one
SELECT r.name as role_name
FROM user_organizations uo
         JOIN roles r ON uo.role_id = r.id
WHERE uo.user_id = ?
  AND uo.org_id = ?;

-- name: CreateInitialRoles :exec
INSERT INTO roles (name, org_id)
VALUES ('user', @org_id),
       ('admin', @org_id),
       ('moderator', @org_id);

-- name: GetModeratorRole :one
SELECT id FROM roles WHERE name = 'moderator' AND org_id = ?;

-- name: GetAdminRole :one
SELECT id FROM roles WHERE name = 'admin' AND org_id = ?;

-- name: GetUserRoleId :one
SELECT id FROM roles WHERE name = 'user' AND org_id = ?;