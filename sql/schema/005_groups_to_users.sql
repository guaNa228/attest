-- +goose Up
ALTER TABLE users
ADD COLUMN "group_id" UUID REFERENCES groups(id) ON DELETE
SET NULL;
-- +goose Down
ALTER TABLE users DROP COLUMN "group_id";