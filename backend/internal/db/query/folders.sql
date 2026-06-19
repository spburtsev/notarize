-- name: GetFolder :one
SELECT * FROM folders
WHERE id = $1;
