-- +goose Up

ALTER TABLE batch_material_usage
DROP CONSTRAINT batch_material_usage_material_lot_fkey;

ALTER TABLE materials
ALTER COLUMN material_lot TYPE TEXT;

ALTER TABLE batch_material_usage
ALTER COLUMN material_lot TYPE TEXT;

ALTER TABLE batch_material_usage
ADD CONSTRAINT batch_material_usage_material_lot_fkey
FOREIGN KEY (material_lot) REFERENCES materials(material_lot);

-- +goose Down

ALTER TABLE batch_material_usage
DROP CONSTRAINT batch_material_usage_material_lot_fkey;

ALTER TABLE batch_material_usage
ALTER COLUMN material_lot TYPE UUID;

ALTER TABLE materials
ALTER COLUMN material_lot TYPE UUID;

ALTER TABLE batch_material_usage
ADD CONSTRAINT batch_material_usage_material_lot_fkey
FOREIGN KEY (material_lot) REFERENCES materials(material_lot);