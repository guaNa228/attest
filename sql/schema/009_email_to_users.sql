-- +goose Up
ALTER TABLE users
ADD COLUMN "email" ON DELETE
SET NULL;
-- +goose Down
ALTER TABLE users DROP COLUMN "email";