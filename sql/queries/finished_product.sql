-- name: AddFinishedProduct :one

INSERT INTO finished_product (
    finished_product_batch,
    In_process_batch_lot,
    component_1_batch,
    component_2_batch,
    created_at,
    updated_at,
    created_by
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;


-- name: GetFinishedProduct :many

SELECT
    finished_product_batch,
    In_process_batch_lot,
    component_1_batch,
    component_2_batch
FROM finished_product;


-- name: GetFinishedProductByIPBatch :one

SELECT
    finished_product_batch,
    In_process_batch_lot,
    component_1_batch,
    component_2_batch
FROM finished_product
WHERE In_process_batch_lot = $1;


-- name: GetFinishedProductByFinishedProductBatch :one

SELECT
    finished_product_batch,
    In_process_batch_lot,
    component_1_batch,
    component_2_batch
FROM finished_product
WHERE finished_product_batch = $1;


-- name: DeleteAllFinishedProducts :exec

DELETE FROM finished_product;


-- name: CheckIPBatchExistsInFill :one

SELECT EXISTS (
    SELECT 1
    FROM fill
    WHERE In_process_batch_lot = $1
);