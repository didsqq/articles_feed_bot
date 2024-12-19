-- +goose Up
-- +goose StatementBegin
ALTER TABLE articles
ADD CONSTRAINT fk_source
FOREIGN KEY (source_id)
REFERENCES sources(id)
ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE articles
DROP CONSTRAINT fk_source;
-- +goose StatementEnd
