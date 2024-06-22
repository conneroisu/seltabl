
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
INSERT INTO
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

/*
 ** selectors.sql
 ** Description: This file contains the SQLite queries for the selectors table
 ** Dialect: sqlite3
 */
/******************************************************************************/
-- name: InsertSelector :one
INSERT INTO
	selectors (value, url_id, context)
VALUES
	(?, ?, ?) RETURNING *;

-- name: UpdateSelectorByID :exec
UPDATE
	selectors
SET
	value = ?,
	url_id = ?,
	context = ?
WHERE
	id = ?;

-- name: DeleteSelectorByID :exec
DELETE FROM
	selectors
WHERE
	id = ?;

-- name: GetSelectorByID :one
SELECT
	*
FROM
	selectors
WHERE
	id = ?;

-- name: GetSelectorByValue :one
SELECT
	*
FROM
	selectors
WHERE
	value = ?;

-- name: GetSelectorsByURL :many
SELECT
	selectors.*
FROM
	selectors
	JOIN urls ON urls.id = selectors.url_id
WHERE
	urls.value = ?;

-- name: GetSelectorsByContext :many
SELECT
	*
FROM
	selectors
WHERE
	context = ?;

/******************************************************************************/

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
