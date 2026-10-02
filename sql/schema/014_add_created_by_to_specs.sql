-- +goose Up

ALTER TABLE specs
ADD COLUMN created_by UUID NOT NULL REFERENCES users(id);

-- +goose Down

ALTER TABLE specs
DROP COLUMN created_by;