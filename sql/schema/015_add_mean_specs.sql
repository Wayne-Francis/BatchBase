-- +goose Up

ALTER TABLE specs
ADD COLUMN mean_min NUMERIC(5,2),
ADD COLUMN mean_max NUMERIC(5,2);

-- +goose Down

ALTER TABLE specs
DROP COLUMN mean_min,
DROP COLUMN mean_max;