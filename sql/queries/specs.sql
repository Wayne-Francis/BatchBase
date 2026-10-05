-- name: AddSpecs :one

INSERT INTO specs (
    test_name,
    created_at,
    updated_at,
    min_result,
    max_result,
    mean_min,
    mean_max,
    rsd_limit,
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
    $8,
    $9
)
RETURNING *;

-- name: GetSpecs :many

SELECT
    test_name,
    created_at,
    updated_at,
    min_result,
    max_result,
    mean_min,
    mean_max,
    rsd_limit,
    created_by
FROM specs;

-- name: GetSpecsByTestName :one

SELECT
    test_name,
    created_at,
    updated_at,
    min_result,
    max_result,
    mean_min,
    mean_max,
    rsd_limit,
    created_by
FROM specs
WHERE test_name = $1;

-- name: DeleteAllSpecs :exec
DELETE FROM specs;

-- name: DeleteSpec :exec
DELETE FROM specs
WHERE test_name = $1;
