/*
** File: files.sql
** Description: This file contains the SQLite queries for the files table
** Dialect: sqlite3
*/
/******************************************************************************/

-- name: ListFiles :many
SELECT
    *
from
    files;

-- name: InsertFile :one
INSERT INTO
    files (uri)
VALUES
    (?) RETURNING *;

-- name: UpdateFileByID :one
UPDATE
    files
SET
    uri = ?
WHERE
    id = ? RETURNING *;

-- name: DeleteFileByID :exec
DELETE FROM
    files
WHERE
    id = ?;

-- name: GetFileByID :one
SELECT
    *
FROM
    files
WHERE
    id = ? LIMIT 1;

-- name: GetFileByURI :one
SELECT
    *
FROM
    files
WHERE
    uri = ?;

/******************************************************************************/
