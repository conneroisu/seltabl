/*
** File: structs.sql
** Description: This file contains the SQLite queries for the structs table
** Dialect: sqlite3
*/
/******************************************************************************/

-- name: ListStructs :many
SELECT
    *
from
    structs;

-- name: InsertStruct :one
INSERT INTO
    structs (file_id, start_line_id, end_line_id, value)
VALUES
    (?, ?, ?, ?) RETURNING *;

-- name: UpdateStructByID :one
UPDATE
    structs
SET
    value = ?,
    start_line_id = ?,
    end_line_id = ?
WHERE
    id = ? RETURNING *;

-- name: DeleteStructByID :exec
DELETE FROM
    structs
WHERE
    id = ?;

-- name: GetStructByID :one
SELECT
    *
FROM
    structs
WHERE
    id = ?; 

-- name: GetStructByValue :one
SELECT
    *
FROM
    structs
WHERE
    value = ?;

-- name: GetStructsByFileID :many
SELECT
    *
FROM
    structs
WHERE
    file_id = ?;

-- name: GetStructsByStartLineID :many
SELECT
    *
FROM
    structs
WHERE
    start_line_id = ?;

-- name: GetStructsByEndLineID :many
SELECT
    *
FROM
    structs
WHERE
    end_line_id = ?;

-- name: GetStructsByValue :many
SELECT
    *
FROM
    structs
WHERE
    value = ?;

-- name: GetStructsByFileIDAndStartLineID :many
SELECT
    *
FROM
    structs
WHERE
    file_id = ? AND start_line_id = ?;

-- name: GetStructsByFileIDAndEndLineID :many
SELECT
    *
FROM
    structs
WHERE
    file_id = ? AND end_line_id = ?;


-- name: GetStructsByStartLineEndlineRange :many
SELECT
    *
FROM
    structs
JOIN lines ON lines.id = structs.start_line_id
JOIN lines AS lines2 ON lines2.id = structs.end_line_id
WHERE
    lines.number <= ? AND lines2.number >= ?;

/******************************************************************************/
