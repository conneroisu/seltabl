
-- name: InsertSelector :one
INSERT INTO selectors (selector, url, context)
VALUES (?, ?, ?) RETURNING *;

-- name: UpdateSelector :exec
UPDATE selectors
SET selector = ?, url = ?, context = ?
WHERE id = ?;

-- name: UpdateSelectorBySelector :exec
UPDATE selectors
SET selector = ?, url = ?, context = ?
WHERE selector = ?;

-- name: ListSelectorsByURL :many
SELECT * from selectors WHERE url = ?;

-- name: DeleteSelector :exec
DELETE FROM selectors
WHERE id = ?;
