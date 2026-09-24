-- +goose Up
CREATE TABLE qc_release (
    finished_product_batch UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    test_name TEXT NOT NULL,
    replicate INTEGER NOT NULL,
    test_date TIMESTAMP NOT NULL,
    result NUMERIC(5,2),

    Primary Key (finished_product_batch, test_name, replicate),
    FOREIGN KEY (finished_product_batch) REFERENCES assembly(finished_product_batch)

);

-- +goose Down
DROP TABLE qc_release;