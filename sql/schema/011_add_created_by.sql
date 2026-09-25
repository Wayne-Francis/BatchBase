-- +goose Up

ALTER TABLE materials
ADD COLUMN created_by UUID NOT NULL REFERENCES users(id);

ALTER TABLE batch_material_usage
ADD COLUMN created_by UUID NOT NULL REFERENCES users(id);

ALTER TABLE blend
ADD COLUMN created_by UUID NOT NULL REFERENCES users(id);

ALTER TABLE fill
ADD COLUMN created_by UUID NOT NULL REFERENCES users(id);

ALTER TABLE assembly
ADD COLUMN created_by UUID NOT NULL REFERENCES users(id);

ALTER TABLE ipc_qc_results
ADD COLUMN created_by UUID NOT NULL REFERENCES users(id);

ALTER TABLE qc_release
ADD COLUMN created_by UUID NOT NULL REFERENCES users(id);

-- +goose Down

ALTER TABLE materials
DROP COLUMN created_by;

ALTER TABLE batch_material_usage
DROP COLUMN created_by;

ALTER TABLE blend
DROP COLUMN created_by;

ALTER TABLE fill
DROP COLUMN created_by;

ALTER TABLE assembly
DROP COLUMN created_by;

ALTER TABLE ipc_qc_results
DROP COLUMN created_by;

ALTER TABLE qc_release
DROP COLUMN created_by;