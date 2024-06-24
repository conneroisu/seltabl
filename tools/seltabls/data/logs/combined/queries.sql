
-- name: GetLogLevelByName :one
SELECT id, name
FROM log_levels
WHERE name = ?;

-- name: InsertLogLevel :one
INSERT INTO log_levels (name)
VALUES (?) RETURNING *;

-- name: UpdateLogLevelById :one
UPDATE log_levels
SET name = ?
WHERE id = ?
RETURNING *;

-- name: DeleteLogLevelById :exec
DELETE FROM log_levels
WHERE id = ?;

/******************************************************************************/

-- name: GetLogById :one
SELECT id, value, created_at, request_id, response_id, notification_id
FROM logs
WHERE id = ?;

-- name: GetLogsByRequestId :many
SELECT id, value, created_at, request_id, response_id, notification_id
FROM logs
WHERE request_id = ?;

-- name: GetLogsByResponseId :many
SELECT id, value, created_at, request_id, response_id, notification_id
FROM logs
WHERE response_id = ?;

-- name: GetLogsByNotificationId :many
SELECT id, value, created_at, request_id, response_id, notification_id
FROM logs
WHERE notification_id = ?;

-- name: InsertLog :one
INSERT INTO logs (value, request_id, response_id, notification_id)
VALUES (?, ?, ?, ?) RETURNING *;

-- name: UpdateLogById :one
UPDATE logs
SET value = ?, request_id = ?, response_id = ?, notification_id = ?
WHERE id = ? RETURNING *;

-- name: DeleteLogById :exec
DELETE FROM logs
WHERE id = ?;

/******************************************************************************/

-- name: GetNotificationById :one
SELECT id, method, created_at, updated_at
FROM notifications
WHERE id = ?;

-- name: GetNotificationsByMethod :many
SELECT id, method, created_at, updated_at
FROM notifications
WHERE method = ?;

-- name: InsertNotification :one
INSERT INTO notifications (method)
VALUES (?) RETURNING *;

-- name: UpdateNotificationById :one
UPDATE notifications
SET method = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? RETURNING *;

-- name: DeleteNotificationById :exec
DELETE FROM notifications
WHERE id = ?;

/******************************************************************************/

-- name: GetRequestById :one
SELECT id, rpc_method, rpc_id
FROM requests
WHERE id = ?;

-- name: GetRequestsByMethod :many
SELECT id, rpc_method, rpc_id
FROM requests
WHERE rpc_method = ?;

-- name: InsertRequest :one
INSERT INTO requests (rpc_method, rpc_id)
VALUES (?, ?) RETURNING *;

-- name: UpdateRequestById :one
UPDATE requests
SET rpc_method = ?, rpc_id = ?
WHERE id = ? RETURNING *;

-- name: DeleteRequestById :exec
DELETE FROM requests
WHERE id = ?;

/******************************************************************************/

-- name: GetResponseById :one
SELECT id, rpc, rpc_id, result, error, created_at
FROM responses
WHERE id = ?;

-- name: GetResponsesByRpcId :many
SELECT id, rpc, rpc_id, result, error, created_at
FROM responses
WHERE rpc_id = ?;

-- name: InsertResponse :one
INSERT INTO responses (rpc, rpc_id, result, error)
VALUES (?, ?, ?, ?) RETURNING *;

-- name: UpdateResponseById :exec
UPDATE responses
SET rpc = ?, rpc_id = ?, result = ?, error = ?
WHERE id = ?;

-- name: DeleteResponseById :exec
DELETE FROM responses
WHERE id = ?;
