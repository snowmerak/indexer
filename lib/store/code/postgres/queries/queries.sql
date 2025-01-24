-- name: GetData :one
SELECT * FROM data WHERE id = $1;

-- name: GetDataList :many
SELECT * FROM data WHERE id = ANY($1::int[]);

-- name: CreateData :one
INSERT INTO data (code_block, file_path, line, description) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateData :one
UPDATE data SET code_block = $2, file_path = $3, line = $4, description = $5 WHERE id = $1 RETURNING *;

-- name: DeleteData :one
DELETE FROM data WHERE id = $1 RETURNING *;

-- name: DeleteDataList :many
DELETE FROM data WHERE id = ANY($1::int[]) RETURNING *;
