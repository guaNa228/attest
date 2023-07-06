-- +goose Up
ALTER TABLE users
ADD COLUMN "email" varchar(255) UNIQUE;
ALTER TABLE users
ADD CONSTRAINT email_format_check CHECK (
        email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
    );
-- +goose Down
ALTER TABLE users DROP CONSTRAINT email_format_check;
ALTER TABLE users DROP COLUMN "email";