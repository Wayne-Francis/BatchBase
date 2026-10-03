-- name: AddFill :one
INSERT INTO fill (
    In_process_batch_lot,
    created_at,
    updated_at,
    fill_start_date,
    fill_end_date,
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


-- name: GetFillList :many
SELECT
    In_process_batch_lot,
    created_at,
    updated_at,
    fill_start_date,
    fill_end_date,
    in_process_disc_hold_end_date,
    created_by
FROM fill;


-- name: GetFillByIPBatch :one
SELECT
    In_process_batch_lot,
    created_at,
    updated_at,
    fill_start_date,
    fill_end_date,
    in_process_disc_hold_end_date,
    created_by
FROM fill
WHERE In_process_batch_lot = $1;


-- name: DeleteAllFill :exec
DELETE FROM fill;

-- name: DeleteIPFromFill :exec
DELETE FROM fill
WHERE In_process_batch_lot = $1;

-- name: CheckIPBatchExistsInBlend :one
SELECT EXISTS (
    SELECT 1
    FROM blend
    WHERE In_process_batch_lot = $1
);

-- name: CheckIPBatchExistsInFinishedProducts :one
SELECT EXISTS (
    SELECT 1
    FROM finished_product
    WHERE In_process_batch_lot = $1
);