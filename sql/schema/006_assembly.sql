-- +goose Up
CREATE TABLE assembly (
    In_process_batch_lot UUID NOT NULL,
    finished_product_batch UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    assembly_start_date TIMESTAMP NOT NULL,
    assembly_end_date TIMESTAMP NOT NULL,
    finished_product_expiry TIMESTAMP GENERATED ALWAYS AS (
    assembly_end_date + INTERVAL '730 days'
    ) STORED,

    FOREIGN KEY (In_process_batch_lot) REFERENCES blend(In_process_batch_lot)

);

-- +goose Down
DROP TABLE assembly;