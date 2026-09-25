-- name: AddRawMaterial :one
INSERT INTO materials (material_lot, created_at, updated_at, material_type, mfg_date, exp_date, created_by)
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

-- name: GetRawMaterials :many
SELECT material_lot, created_at, updated_at, material_type, mfg_date, exp_date, created_by
FROM materials;

-- name: GetRawMaterialByLot :one
SELECT *
FROM materials
WHERE material_lot = $1;

-- name: DeleteMaterials :exec
DELETE FROM materials;
