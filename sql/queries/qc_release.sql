-- name: AddQCReleaseResults :one

INSERT INTO qc_release (
    finished_product_batch,
    created_at,
    updated_at,
    test_name,
    replicate,
    test_date,
    result,
    created_by
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;


-- name: GetAllQCReleaseResults :many

SELECT
    finished_product_batch,
    created_at,
    updated_at,
    test_name,
    replicate,
    test_date,
    result
FROM qc_release;


-- name: GetQCReleaseResultsByFPBatch :many

SELECT
    finished_product_batch,
    created_at,
    updated_at,
    test_name,
    replicate,
    test_date,
    result
FROM qc_release
WHERE finished_product_batch = $1;

-- name: CheckFPBatchExists :one

SELECT EXISTS (
    SELECT 1
    FROM assembly
    WHERE finished_product_batch = $1
);

-- name: DeleteAllQCReleaseResults :exec

DELETE FROM qc_release;

-- name: GetQCReleaseResultsByIPBatch :many

SELECT
    qr.finished_product_batch,
    fp.In_process_batch_lot,
    qr.test_name,
    qr.replicate,
    qr.test_date,
    qr.result
FROM qc_release qr
JOIN finished_product fp
    ON fp.finished_product_batch = qr.finished_product_batch
WHERE fp.In_process_batch_lot = $1;

