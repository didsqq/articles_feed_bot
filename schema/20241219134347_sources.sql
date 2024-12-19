-- +goose Up
-- +goose StatementBegin
ALTER TABLE sources DROP COLUMN updated_at;
ALTER TABLE sources ADD COLUMN priority INT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sources ADD COLUMN updated_at NOT NULL DEFAULT NOW();
ALTER TABLE sources DROP COLUMN priority;
-- +goose StatementEnd
