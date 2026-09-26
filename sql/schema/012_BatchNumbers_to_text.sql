-- +goose Up

-- Drop foreign keys that depend on the batch-number columns

ALTER TABLE batch_material_usage
DROP CONSTRAINT batch_material_usage_material_lot_fkey;

ALTER TABLE fill
DROP CONSTRAINT fill_In_process_batch_lot_fkey;

ALTER TABLE assembly
DROP CONSTRAINT assembly_In_process_batch_lot_fkey;

ALTER TABLE ipc_qc_results
DROP CONSTRAINT ipc_qc_results_In_process_batch_lot_fkey;

ALTER TABLE qc_release
DROP CONSTRAINT qc_release_finished_product_batch_fkey;


-- Change In-Process Batch Lot from UUID to TEXT

ALTER TABLE batch_material_usage
ALTER COLUMN In_process_batch_lot TYPE TEXT
USING In_process_batch_lot::TEXT;

ALTER TABLE blend
ALTER COLUMN In_process_batch_lot TYPE TEXT
USING In_process_batch_lot::TEXT;

ALTER TABLE fill
ALTER COLUMN In_process_batch_lot TYPE TEXT
USING In_process_batch_lot::TEXT;

ALTER TABLE assembly
ALTER COLUMN In_process_batch_lot TYPE TEXT
USING In_process_batch_lot::TEXT;

ALTER TABLE ipc_qc_results
ALTER COLUMN In_process_batch_lot TYPE TEXT
USING In_process_batch_lot::TEXT;


-- Change Finished Product Batch from UUID to TEXT

ALTER TABLE assembly
ALTER COLUMN finished_product_batch TYPE TEXT
USING finished_product_batch::TEXT;

ALTER TABLE qc_release
ALTER COLUMN finished_product_batch TYPE TEXT
USING finished_product_batch::TEXT;


-- Re-create the foreign keys

ALTER TABLE batch_material_usage
ADD CONSTRAINT batch_material_usage_material_lot_fkey
FOREIGN KEY (material_lot) REFERENCES materials(material_lot);

ALTER TABLE fill
ADD CONSTRAINT fill_In_process_batch_lot_fkey
FOREIGN KEY (In_process_batch_lot) REFERENCES blend(In_process_batch_lot);

ALTER TABLE assembly
ADD CONSTRAINT assembly_In_process_batch_lot_fkey
FOREIGN KEY (In_process_batch_lot) REFERENCES blend(In_process_batch_lot);

ALTER TABLE ipc_qc_results
ADD CONSTRAINT ipc_qc_results_In_process_batch_lot_fkey
FOREIGN KEY (In_process_batch_lot) REFERENCES blend(In_process_batch_lot);

ALTER TABLE qc_release
ADD CONSTRAINT qc_release_finished_product_batch_fkey
FOREIGN KEY (finished_product_batch) REFERENCES assembly(finished_product_batch);


-- +goose Down

-- Drop foreign keys

ALTER TABLE fill
DROP CONSTRAINT fill_In_process_batch_lot_fkey;

ALTER TABLE assembly
DROP CONSTRAINT assembly_In_process_batch_lot_fkey;

ALTER TABLE ipc_qc_results
DROP CONSTRAINT ipc_qc_results_In_process_batch_lot_fkey;

ALTER TABLE qc_release
DROP CONSTRAINT qc_release_finished_product_batch_fkey;


-- Change Finished Product Batch back to UUID

ALTER TABLE qc_release
ALTER COLUMN finished_product_batch TYPE UUID
USING finished_product_batch::UUID;

ALTER TABLE assembly
ALTER COLUMN finished_product_batch TYPE UUID
USING finished_product_batch::UUID;


-- Change In-Process Batch Lot back to UUID

ALTER TABLE ipc_qc_results
ALTER COLUMN In_process_batch_lot TYPE UUID
USING In_process_batch_lot::UUID;

ALTER TABLE assembly
ALTER COLUMN In_process_batch_lot TYPE UUID
USING In_process_batch_lot::UUID;

ALTER TABLE fill
ALTER COLUMN In_process_batch_lot TYPE UUID
USING In_process_batch_lot::UUID;

ALTER TABLE blend
ALTER COLUMN In_process_batch_lot TYPE UUID
USING In_process_batch_lot::UUID;

ALTER TABLE batch_material_usage
ALTER COLUMN In_process_batch_lot TYPE UUID
USING In_process_batch_lot::UUID;


-- Re-create foreign keys

ALTER TABLE fill
ADD CONSTRAINT fill_In_process_batch_lot_fkey
FOREIGN KEY (In_process_batch_lot) REFERENCES blend(In_process_batch_lot);

ALTER TABLE assembly
ADD CONSTRAINT assembly_In_process_batch_lot_fkey
FOREIGN KEY (In_process_batch_lot) REFERENCES blend(In_process_batch_lot);

ALTER TABLE ipc_qc_results
ADD CONSTRAINT ipc_qc_results_In_process_batch_lot_fkey
FOREIGN KEY (In_process_batch_lot) REFERENCES blend(In_process_batch_lot);

ALTER TABLE qc_release
ADD CONSTRAINT qc_release_finished_product_batch_fkey
FOREIGN KEY (finished_product_batch) REFERENCES assembly(finished_product_batch);