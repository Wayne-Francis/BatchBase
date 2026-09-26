-- name: AddMaterialUsage :one

INSERT INTO batch_material_usage (
    In_process_batch_lot,
    material_lot,
    created_at,
    updated_at,
    created_by
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;


-- name: GetMaterialUsageList :many

SELECT
    In_process_batch_lot,
    material_lot,
    created_at,
    updated_at,
    created_by
FROM batch_material_usage;


-- name: GetMaterialUsageByIPBatch :many

SELECT *
FROM batch_material_usage
WHERE In_process_batch_lot = $1;


-- name: GetIPBatchesByMaterialLot :many

SELECT *
FROM batch_material_usage
WHERE material_lot = $1;


-- name: DeleteAllMaterialUsage :exec

DELETE FROM batch_material_usage;