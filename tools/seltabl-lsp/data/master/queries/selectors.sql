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
