-- +goose Up
CREATE TABLE materials (
    material_lot UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    material_type TEXT NOT NULL,
    mfg_date TIMESTAMP NOT NULL,
    exp_date TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE materials;