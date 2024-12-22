-- name: CreateAccessRequest :one
INSERT INTO access_requests (user_id, resource_id, status, reason, expires_at)
VALUES (?,
        ?,
        ?,
        ?,
        ?)
RETURNING *;

-- name: GetAccessRequest :one
SELECT *
FROM access_requests
WHERE id = ?;

-- name: UpdateAccessRequestStatus :exec
UPDATE access_requests
SET status = ?
WHERE id = ?;

-- name: ListUserAccessRequests :many
SELECT *
FROM access_requests
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: ListPendingAccessRequests :many
SELECT ar.*, u.email, r.name as resource_name
FROM access_requests ar
         JOIN users u ON ar.user_id = u.id
         JOIN resources r ON ar.resource_id = r.id
WHERE ar.status = 'pending' AND r.org_id = ?
ORDER BY ar.created_at DESC;

-- name: ListActiveUserAccess :many
SELECT
    ar.*,
    r.name as resource_name
FROM access_requests ar
         JOIN resources r ON ar.resource_id = r.id
WHERE ar.user_id = ?
  AND ar.expires_at > ?
  AND ar.status = 'approved'
ORDER BY ar.expires_at;

-- name: RevokeExpiredAccess :exec
UPDATE access_requests
SET status = 'expired'
WHERE expires_at < current_time
  AND status = 'approved';

-- name: GetActiveAccessRequest :one
SELECT * FROM access_requests
WHERE user_id = ?
  AND resource_id = ?
  AND expires_at > ?
  AND status = 'approved';