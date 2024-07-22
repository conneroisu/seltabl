
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

/*
 ** selectors.sql
 ** Description: This file contains the SQLite queries for the selectors table
 ** Dialect: sqlite3
 */
/******************************************************************************/
-- name: InsertSelector :one
INSERT INTO
	selectors (value, url_id, context, occurances)
VALUES
	(?, ?, ?, ?) RETURNING *;

-- name: UpdateSelectorByID :exec
UPDATE
	selectors
SET
	value = ?,
	url_id = ?,
	context = ?,
	occurances = ?
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

-- name: GetSelectorsByMinOccurances :many
SELECT
    selectors.*
FROM
    selectors
    JOIN urls ON urls.id = selectors.url_id
WHERE
    selectors.occurances >= ? AND
    urls.value = ?;
/******************************************************************************/

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
INSERT OR IGNORE INTO
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
