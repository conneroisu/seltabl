/*
** File: tags.sql
** Description: This file contains the SQLite queries for the tags table
** Dialect: sqlite3
*/
/******************************************************************************/

-- name: ListTags :many
SELECT
    *
from
    tags;

-- name: InsertTag :exec
INSERT INTO
    tags (value, start, end, line_id, field_id)
VALUES
    (?, ?, ?, ?, ?);

-- name: UpdateTagByID :exec
UPDATE
    tags
SET
    value = ?,
    start = ?,
    end = ?,
    line_id = ?,
    field_id = ?
WHERE
    id = ?;

-- name: DeleteTagByID :exec
DELETE FROM
    tags
WHERE
    id = ?;

-- name: GetTagByID :one
SELECT
    *
FROM
    tags
WHERE
    id = ?;

-- name: GetTagByValue :one
SELECT
    *
FROM
    tags
WHERE
    value = ?;

-- name: GetTagsByFieldID :many
SELECT
    *
FROM
    tags
WHERE
    field_id = ?;

-- name: GetTagByFieldIDAndValue :one
SELECT
    *
FROM
    tags
WHERE
    field_id = ? AND value = ?;
