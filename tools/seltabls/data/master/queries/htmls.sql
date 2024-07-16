/* 
 ** File: htmls.sql
 ** Description: This file contains the SQLite queries for the htmls table
 ** Dialect: sqlite3
 */
/******************************************************************************/
-- name: ListHTMLs :many
SELECT
    *
from
    htmls;

-- name: InsertHTML :one
INSERT OR IGNORE INTO
    htmls (value)
VALUES
    (?) RETURNING *;

-- name: UpdateHTMLByID :one
UPDATE
    htmls
SET
    value = ?
WHERE
    id = ? RETURNING *;

-- name: DeleteHTMLByID :exec
DELETE FROM
    htmls
WHERE
    id = ?;

/******************************************************************************/
