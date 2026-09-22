-- +goose Up
CREATE TABLE fill (
    In_process_batch_lot UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    fill_start_date TIMESTAMP NOT NULL,
    fill_end_date TIMESTAMP NOT NULL,
    in_process_disc_hold_end_date TIMESTAMP GENERATED ALWAYS AS (
    fill_end_date + INTERVAL '90 days'
    ) STORED,

    FOREIGN KEY (In_process_batch_lot) REFERENCES blend(In_process_batch_lot)

);

-- +goose Down
DROP TABLE fill;