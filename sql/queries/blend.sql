-- name: AddBlend :one
INSERT INTO blend (
    In_process_batch_lot,
    created_at,
    updated_at,
    blend_start_date,
    blend_end_date,
    created_by
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;


-- name: GetBlendList :many
SELECT
    In_process_batch_lot,
    created_at,
    updated_at,
    blend_start_date,
    blend_end_date,
    in_process_blend_hold_end_date,
    created_by
FROM blend;


-- name: GetBlendByIPBatch :one
SELECT
    In_process_batch_lot,
    created_at,
    updated_at,
    blend_start_date,
    blend_end_date,
    in_process_blend_hold_end_date,
    created_by
FROM blend
WHERE In_process_batch_lot = $1;


-- name: DeleteAllBlend :exec
DELETE FROM blend;

-- name: CheckIPBatchExists :one
SELECT EXISTS (
    SELECT 1
    FROM batch_material_usage
    WHERE In_process_batch_lot = $1
);