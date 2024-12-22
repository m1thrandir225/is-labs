-- name: CreateRolePermission :one
INSERT INTO role_permissions (role_id, resource_id, can_read, can_write, can_delete)
VALUES (?,
        ?,
        ?,
        ?,
        ?)
RETURNING *;

-- name: GetRolePermissions :one
SELECT *
FROM role_permissions
WHERE role_id = ?
  AND resource_id = ?;

-- name: UpdateRolePermissions :one
UPDATE role_permissions
SET can_read   = ?,
    can_write  = ?,
    can_delete = ?
WHERE role_id = ?
  AND resource_id = ?
RETURNING *;