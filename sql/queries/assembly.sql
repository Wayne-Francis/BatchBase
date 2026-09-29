-- name: AddAssembly :one

INSERT INTO assembly (
    finished_product_batch,
    created_at,
    updated_at,
    assembly_start_date,
    assembly_end_date,
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


-- name: GetAssembly :many

SELECT
    finished_product_batch,
    assembly_start_date,
    assembly_end_date,
    finished_product_expiry
FROM assembly;


-- name: GetAssemblyByIPBatch :one

SELECT
    a.finished_product_batch,
    fp.In_process_batch_lot,
    a.assembly_start_date,
    a.assembly_end_date,
    a.finished_product_expiry
FROM assembly a
JOIN finished_product fp
    ON a.finished_product_batch = fp.finished_product_batch
WHERE fp.In_process_batch_lot = $1;


-- name: GetAssemblyByFinishedProductBatch :one

SELECT
    finished_product_batch,
    assembly_start_date,
    assembly_end_date,
    finished_product_expiry
FROM assembly
WHERE finished_product_batch = $1;

-- name: CheckFPBatchExistsInFinishedProduct :one

SELECT EXISTS (
    SELECT 1
    FROM finished_product
    WHERE finished_product_batch = $1
);

-- name: DeleteAllAssembly :exec

DELETE FROM assembly;

