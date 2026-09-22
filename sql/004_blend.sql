-- +goose Up
CREATE TABLE blend (
    In_process_batch_lot UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    blend_start_date TIMESTAMP NOT NULL,
    blend_end_date TIMESTAMP NOT NULL,
    in_process_blend_hold_end_date TIMESTAMP GENERATED ALWAYS AS (
    blend_end_date + INTERVAL '7 days'
    ) STORED

);

-- +goose Down
DROP TABLE blend;