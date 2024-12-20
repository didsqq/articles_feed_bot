-- +goose Up
-- +goose StatementBegin
ALTER TABLE articles ADD CONSTRAINT unique_link UNIQUE (link);
ALTER TABLE articles ALTER COLUMN created_at SET DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE articles DROP CONSTRAINT unique_link;
ALTER TABLE articles ALTER COLUMN created_at DROP DEFAULT NOW();
-- +goose StatementEnd
