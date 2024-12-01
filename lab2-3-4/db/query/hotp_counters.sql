

-- name: GetCurrentCounter :one
SELECT counter
FROM hotp_counters
WHERE user_id = ?;

-- name: CreateHotpCounter :exec
INSERT INTO hotp_counters (user_id, counter)
VALUES (?, ?)
ON CONFLICT (user_id) DO NOTHING;

-- name: CleanupExpiredCounters :exec
DELETE FROM hotp_counters
WHERE last_used_timestamp < datetime('now', '-30 days');

-- name: IncreaseCounter :one
UPDATE hotp_counters
SET counter = counter + 1, last_used_timestamp = CURRENT_TIMESTAMP
WHERE user_id = ?
RETURNING  counter;