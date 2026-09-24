-- +goose Up
CREATE TABLE specs (
    test_name TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    min_result NUMERIC(5,2),
    max_result NUMERIC(5,2),
    rsd_limit NUMERIC(5,2)

);

-- +goose Down
DROP TABLE specs;