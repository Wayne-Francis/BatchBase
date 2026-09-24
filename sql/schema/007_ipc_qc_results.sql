-- +goose Up
CREATE TABLE ipc_qc_results (
    In_process_batch_lot UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    test_name TEXT NOT NULL,
    replicate INTEGER NOT NULL,
    test_date TIMESTAMP NOT NULL,
    result NUMERIC(5,2),

    Primary Key (In_process_batch_lot, test_name, replicate),
    FOREIGN KEY (In_process_batch_lot) REFERENCES blend(In_process_batch_lot)

);

-- +goose Down
DROP TABLE ipc_qc_results;