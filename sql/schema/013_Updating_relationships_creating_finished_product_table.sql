-- +goose Up

-- Drop the existing QC Release -> Assembly foreign key
ALTER TABLE qc_release
DROP CONSTRAINT qc_release_finished_product_batch_fkey;


-- Create Finished Product
CREATE TABLE finished_product (
    finished_product_batch TEXT PRIMARY KEY,
    In_process_batch_lot TEXT NOT NULL,
    component_1_batch TEXT NOT NULL,
    component_2_batch TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_by UUID NOT NULL,   

    FOREIGN KEY (In_process_batch_lot)
        REFERENCES blend(In_process_batch_lot),

    FOREIGN KEY (created_by)
        REFERENCES users(id)
);


-- Remove IP batch from Assembly
ALTER TABLE assembly
DROP CONSTRAINT assembly_In_process_batch_lot_fkey;

ALTER TABLE assembly
DROP COLUMN In_process_batch_lot;


-- Link Assembly to Finished Product
ALTER TABLE assembly
ADD CONSTRAINT assembly_finished_product_batch_fkey
FOREIGN KEY (finished_product_batch)
REFERENCES finished_product(finished_product_batch);


-- Link QC Release to Finished Product
ALTER TABLE qc_release
ADD CONSTRAINT qc_release_finished_product_batch_fkey
FOREIGN KEY (finished_product_batch)
REFERENCES finished_product(finished_product_batch);


-- +goose Down

-- Remove QC Release -> Finished Product
ALTER TABLE qc_release
DROP CONSTRAINT qc_release_finished_product_batch_fkey;


-- Remove Assembly -> Finished Product
ALTER TABLE assembly
DROP CONSTRAINT assembly_finished_product_batch_fkey;


-- Restore Assembly -> Blend relationship
ALTER TABLE assembly
ADD COLUMN In_process_batch_lot TEXT NOT NULL;

ALTER TABLE assembly
ADD CONSTRAINT assembly_In_process_batch_lot_fkey
FOREIGN KEY (In_process_batch_lot)
REFERENCES blend(In_process_batch_lot);


-- Restore QC Release -> Assembly relationship
ALTER TABLE qc_release
ADD CONSTRAINT qc_release_finished_product_batch_fkey
FOREIGN KEY (finished_product_batch)
REFERENCES assembly(finished_product_batch);


-- Remove Finished Product
DROP TABLE finished_product;