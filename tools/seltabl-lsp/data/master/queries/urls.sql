/*
 ** File: urls.sql
 ** Description: This file contains the SQLite queries for the urls table
 ** Dialect: sqlite3
 */
/******************************************************************************/
-- name: ListURLs :many
SELECT
    *
from
    urls;

-- name: GetURLByValue :one
SELECT
    *
FROM
    urls
WHERE
    value = ?;

-- name: InsertURL :one
INSERT INTO
    urls (value, html_id)
VALUES
    (?, ?) RETURNING *;

-- name: UpdateURL :exec
UPDATE
    urls
SET
    value = ?,
    html_id = ?
WHERE
    id = ?;

-- name: DeleteURL :exec
DELETE FROM
    urls
WHERE
    id = ?;

-- name: ListAll :many
SELECT
    urls.id,
    urls.value,
    htmls.value as html,
    selectors.value as selector
FROM
    urls
    JOIN htmls ON urls.html_id = htmls.id
    JOIN selectors ON urls.id = selectors.url_id;

-- name: UpsertURL :one
INSERT INTO
    urls (value, html_id)
VALUES
    (?, ?)
ON CONFLICT (value)
DO UPDATE
    SET
        html_id = excluded.html_id RETURNING *;
/******************************************************************************/
