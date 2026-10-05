-- +goose Up

ALTER TABLE ipc_qc_results
ALTER COLUMN result SET NOT NULL;

ALTER TABLE qc_release
ALTER COLUMN result SET NOT NULL;


-- +goose Down

ALTER TABLE ipc_qc_results
ALTER COLUMN result DROP NOT NULL;

ALTER TABLE qc_release
ALTER COLUMN result DROP NOT NULL;