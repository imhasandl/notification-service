-- name: RegisterDeviceToken :one
INSERT INTO device_tokens(id, user_id, device_token, device_type, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT (user_id, device_token)
DO UPDATE SET updated_at = NOW(), device_type = $3
RETURNING *;

-- name: GetDeviceTokensByUser :one
SELECT device_token FROM device_tokens
WHERE user_id = $1;

-- name: DeleteDeviceToken :exec
DELETE FROM device_tokens
WHERE user_id = $1 AND device_token = $2;