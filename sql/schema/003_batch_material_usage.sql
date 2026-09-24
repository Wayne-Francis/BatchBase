-- +goose Up
CREATE TABLE batch_material_usage (
    In_process_batch_lot UUID NOT NULL,
    material_lot UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    Primary Key (In_process_batch_lot, material_lot),
    FOREIGN KEY (material_lot) REFERENCES materials(material_lot)
);

-- +goose Down
DROP TABLE batch_material_usage;