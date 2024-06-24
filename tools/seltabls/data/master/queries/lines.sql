/*
** File: lines.sql
** Description: This file contains the SQLite queries for the lines table
** Dialect: sqlite3
*/
/******************************************************************************/

-- name: ListLines :many
SELECT
    *
from
    lines;

-- name: InsertLine :one
INSERT INTO
    lines (value)
VALUES
    (?) RETURNING id, value, number;

-- name: UpdateLineByID :one
UPDATE
    lines
SET
    value = ?
WHERE
    id = ? RETURNING id, value, number;

-- name: DeleteLineByID :exec
DELETE FROM
    lines
WHERE
    id = ?;

-- name: GetLineByID :one
SELECT
    id, value, number
FROM
    lines
WHERE
    id = ?;

/******************************************************************************/
