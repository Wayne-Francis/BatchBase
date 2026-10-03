-- name: AddIPCQCResults :one

INSERT INTO ipc_qc_results (
    In_process_batch_lot,
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


-- name: GetIPCQCResults :many

SELECT
    In_process_batch_lot,
    created_at,
    updated_at,
    test_name,
    replicate,
    test_date,
    result
FROM ipc_qc_results;


-- name: GetIPCQCResultsByIPBatch :many

SELECT
    In_process_batch_lot,
    created_at,
    updated_at,
    test_name,
    replicate,
    test_date,
    result
FROM ipc_qc_results
WHERE In_process_batch_lot = $1;

-- name: CheckIPCBatchExists :one

SELECT EXISTS (
    SELECT 1
    FROM blend
    WHERE In_process_batch_lot = $1
);

-- name: DeleteAllIPCQCResults :exec

DELETE FROM ipc_qc_results;

-- name: DeleteIPCQCResultsForIPBatch :exec

DELETE FROM ipc_qc_results
WHERE In_process_batch_lot = $1;

-- name: GetIPCQCResultsByFPBatch :many

SELECT
    ipc.In_process_batch_lot,
    fp.finished_product_batch,
    ipc.test_name,
    ipc.replicate,
    ipc.test_date,
    ipc.result
FROM ipc_qc_results ipc
JOIN finished_product fp
    ON ipc.In_process_batch_lot = fp.In_process_batch_lot
WHERE fp.finished_product_batch = $1;

-- name: CheckIPCBatchExistsInIPCQCResults :one

SELECT EXISTS (
    SELECT 1
    FROM ipc_qc_results
    WHERE In_process_batch_lot = $1
);

-- name: CountIPCQCResults :one
SELECT COUNT(*)
FROM ipc_qc_results
WHERE In_process_batch_lot = $1
AND test_name = $2;