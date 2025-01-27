-- name: GetData :one
SELECT * FROM data WHERE project = $1 AND id = $2;

-- name: GetDataList :many
SELECT * FROM data WHERE project = $1 AND id = ANY($2::int[]);

-- name: CreateData :one
INSERT INTO data (project, id, code_block, file_path, line, description) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateData :one
UPDATE data SET code_block = $3, file_path = $4, line = $5, description = $6 WHERE project = $1 AND id = $2 RETURNING *;

-- name: DeleteData :one
DELETE FROM data WHERE project = $1 AND id = $2 RETURNING *;

-- name: DeleteProjectData :many
DELETE FROM data WHERE project = $1 RETURNING *;

-- name: DeleteDataList :many
DELETE FROM data WHERE project = $1 AND id = ANY($2::int[]) RETURNING *;
